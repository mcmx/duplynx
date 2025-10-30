package actions

import "sync"

// Store maintains duplicate groups in-memory for demo flows.
type Store struct {
	mu     sync.RWMutex
	groups map[string]*DuplicateGroup
}

func NewStore(groups []DuplicateGroup) *Store {
	m := make(map[string]*DuplicateGroup, len(groups))
	for i := range groups {
		g := groups[i]
		tmp := g
		m[g.ID] = &tmp
	}
	return &Store{groups: m}
}

func (s *Store) Get(id string) (*DuplicateGroup, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	g, ok := s.groups[id]
	return g, ok
}

func (s *Store) ListByScan(scanID string) []DuplicateGroup {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []DuplicateGroup
	for _, g := range s.groups {
		if g.ScanID == scanID {
			cp := *g
			out = append(out, cp)
		}
	}
	return out
}

func (s *Store) Update(group *DuplicateGroup) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if existing, ok := s.groups[group.ID]; ok {
		*existing = *group
	}
}
