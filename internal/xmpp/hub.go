package xmpp

import (
	"net/url"
	"strings"
	"sync"
)

type client struct {
	conn        *conn
	accountID   string
	displayName string
	jid         string
	resource    string
}

func (c *client) nick() string {
	return url.PathEscape(c.displayName) + ":" + c.accountID + ":" + c.resource
}

func (c *client) send(e element) {
	_ = c.conn.write(e.String())
}

type hub struct {
	mu      sync.RWMutex
	clients map[string]*client
	rooms   map[string]map[string]*client
}

func newHub() *hub {
	return &hub{
		clients: make(map[string]*client),
		rooms:   make(map[string]map[string]*client),
	}
}

func (h *hub) add(c *client) {
	h.mu.Lock()
	previous := h.clients[c.accountID]
	h.clients[c.accountID] = c
	h.mu.Unlock()

	if previous != nil {
		previous.conn.close()
	}
}

func (h *hub) remove(c *client) []string {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[c.accountID] == c {
		delete(h.clients, c.accountID)
	}
	var left []string
	for name, members := range h.rooms {
		if members[c.accountID] != c {
			continue
		}
		delete(members, c.accountID)
		left = append(left, name)
		if len(members) == 0 {
			delete(h.rooms, name)
		}
	}
	return left
}

func (h *hub) find(jid string) *client {
	account, _, _ := strings.Cut(jid, "@")

	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.clients[account]
}

func (h *hub) join(room string, c *client) ([]*client, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	members := h.rooms[room]
	if members == nil {
		members = make(map[string]*client)
		h.rooms[room] = members
	}
	if members[c.accountID] == c {
		return nil, false
	}
	members[c.accountID] = c
	return snapshot(members), true
}

func (h *hub) leave(room string, c *client) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	members := h.rooms[room]
	if members[c.accountID] != c {
		return false
	}
	delete(members, c.accountID)
	if len(members) == 0 {
		delete(h.rooms, room)
	}
	return true
}

func (h *hub) members(room string) []*client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return snapshot(h.rooms[room])
}

func snapshot(members map[string]*client) []*client {
	list := make([]*client, 0, len(members))
	for _, m := range members {
		list = append(list, m)
	}
	return list
}

func (h *hub) closeAll() {
	h.mu.RLock()
	clients := make([]*client, 0, len(h.clients))
	for _, c := range h.clients {
		clients = append(clients, c)
	}
	h.mu.RUnlock()

	for _, c := range clients {
		c.conn.close()
	}
}

func (h *hub) client(accountID string) *client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.clients[accountID]
}
