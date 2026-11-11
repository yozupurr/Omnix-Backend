package party

import (
	"encoding/json"
	"maps"
	"net/http"
	"sync"
	"time"

	"omnix/internal/platform/httpx"
)

const maxBody = 1 << 16

type store struct {
	mu      sync.Mutex
	parties map[string]*party
}

type userParties struct {
	Current []*party `json:"current"`
	Pending []any    `json:"pending"`
	Invites []any    `json:"invites"`
	Pings   []any    `json:"pings"`
}

type undeliveredCount struct {
	Pings   int `json:"pings"`
	Invites int `json:"invites"`
}

type createRequest struct {
	Config   map[string]any `json:"config"`
	JoinInfo joinInfo       `json:"join_info"`
	Meta     meta           `json:"meta"`
}

type updateRequest struct {
	Config map[string]any `json:"config"`
	Meta   metaPatch      `json:"meta"`
}

type joinResponse struct {
	Status  string `json:"status"`
	PartyID string `json:"party_id"`
}

func Register(mux *http.ServeMux) {
	s := &store{parties: map[string]*party{}}
	mux.HandleFunc("GET /party/api/v1/Fortnite/user/{accountId}", s.user)
	mux.HandleFunc("GET /party/api/v1/Fortnite/user/{accountId}/notifications/undelivered/count", undelivered)
	mux.HandleFunc("POST /party/api/v1/Fortnite/parties", s.create)
	mux.HandleFunc("GET /party/api/v1/Fortnite/parties/{partyId}", s.get)
	mux.HandleFunc("PATCH /party/api/v1/Fortnite/parties/{partyId}", s.update)
	mux.HandleFunc("PATCH /party/api/v1/Fortnite/parties/{partyId}/members/{accountId}/meta", s.updateMember)
	mux.HandleFunc("POST /party/api/v1/Fortnite/parties/{partyId}/members/{accountId}/join", s.join)
	mux.HandleFunc("POST /party/api/v1/Fortnite/parties/{partyId}/members/{accountId}/promote", s.promote)
	mux.HandleFunc("DELETE /party/api/v1/Fortnite/parties/{partyId}/members/{accountId}", s.leave)
}

func (s *store) user(w http.ResponseWriter, r *http.Request) {
	accountID := r.PathValue("accountId")
	s.mu.Lock()
	defer s.mu.Unlock()
	current := []*party{}
	for _, p := range s.parties {
		if _, ok := p.member(accountID); ok {
			current = append(current, p)
		}
	}
	httpx.JSON(w, http.StatusOK, userParties{Current: current, Pending: []any{}, Invites: []any{}, Pings: []any{}})
}

func undelivered(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, undeliveredCount{})
}

func (s *store) create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if !decode(w, r, &req) {
		return
	}
	accountID := req.JoinInfo.accountID()
	if accountID == "" {
		httpx.InvalidRequest.Write(w)
		return
	}
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	p := newParty(req.Config, req.Meta, now)
	s.removeEverywhere(accountID)
	p.join(accountID, req.JoinInfo, now)
	s.parties[p.ID] = p
	httpx.JSON(w, http.StatusOK, p)
}

func (s *store) get(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.find(w, r)
	if !ok {
		return
	}
	httpx.JSON(w, http.StatusOK, p)
}

func (s *store) update(w http.ResponseWriter, r *http.Request) {
	var req updateRequest
	if !decode(w, r, &req) {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.find(w, r)
	if !ok {
		return
	}
	maps.Copy(p.Config, req.Config)
	p.Meta.apply(req.Meta)
	p.touch(time.Now())
	w.WriteHeader(http.StatusNoContent)
}

func (s *store) updateMember(w http.ResponseWriter, r *http.Request) {
	var req metaPatch
	if !decode(w, r, &req) {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.find(w, r)
	if !ok {
		return
	}
	m, ok := p.member(r.PathValue("accountId"))
	if !ok {
		httpx.NotFound.Write(w)
		return
	}
	now := time.Now()
	m.Meta.apply(req)
	m.Revision++
	m.UpdatedAt = httpx.Timestamp(now)
	p.touch(now)
	w.WriteHeader(http.StatusNoContent)
}

func (s *store) join(w http.ResponseWriter, r *http.Request) {
	var req joinInfo
	if !decode(w, r, &req) {
		return
	}
	accountID := r.PathValue("accountId")
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.find(w, r)
	if !ok {
		return
	}
	s.removeEverywhere(accountID)
	p.join(accountID, req, time.Now())
	httpx.JSON(w, http.StatusOK, joinResponse{Status: "JOINED", PartyID: p.ID})
}

func (s *store) promote(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.find(w, r)
	if !ok {
		return
	}
	if !p.promote(r.PathValue("accountId")) {
		httpx.NotFound.Write(w)
		return
	}
	p.touch(time.Now())
	w.WriteHeader(http.StatusNoContent)
}

func (s *store) leave(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.find(w, r)
	if !ok {
		return
	}
	if p.leave(r.PathValue("accountId")) {
		p.touch(time.Now())
	}
	if len(p.Members) == 0 {
		delete(s.parties, p.ID)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *store) find(w http.ResponseWriter, r *http.Request) (*party, bool) {
	p, ok := s.parties[r.PathValue("partyId")]
	if !ok {
		httpx.PartyNotFound.Write(w)
	}
	return p, ok
}

func (s *store) removeEverywhere(accountID string) {
	for id, p := range s.parties {
		if p.leave(accountID) && len(p.Members) == 0 {
			delete(s.parties, id)
		}
	}
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBody)).Decode(v); err != nil {
		httpx.InvalidRequest.Write(w)
		return false
	}
	return true
}
