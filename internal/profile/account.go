package profile

import "sync"

const startingVBucks = 10000

var starterItems = []string{
	"AthenaPickaxe:DefaultPickaxe",
	"AthenaGlider:DefaultGlider",
	"AthenaDance:EID_DanceMoves",
}

type account struct {
	vbucks    int
	athena    map[string]int
	seen      map[string]struct{}
	banners   map[string]int
	favorites map[string]any
	books     map[int]*book
	revisions map[string]int
}

type book struct {
	purchased bool
	level     int
}

func newAccount() *account {
	a := &account{
		vbucks:    startingVBucks,
		athena:    make(map[string]int, len(starterItems)),
		seen:      make(map[string]struct{}, len(starterItems)),
		banners:   map[string]int{},
		favorites: newFavorites(),
		books:     map[int]*book{},
		revisions: map[string]int{},
	}
	for _, id := range starterItems {
		a.athena[id] = 1
		a.seen[id] = struct{}{}
	}
	return a
}

func (a *account) book(season int) *book {
	b, ok := a.books[season]
	if !ok {
		b = &book{level: 1}
		a.books[season] = b
	}
	return b
}

type store struct {
	mu       sync.Mutex
	accounts map[string]*account
}

func newStore() *store {
	return &store{accounts: map[string]*account{}}
}

func (s *store) with(accountID string, fn func(*account)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.accounts[accountID]
	if !ok {
		a = newAccount()
		s.accounts[accountID] = a
	}
	fn(a)
}
