package store

import (
	"coupon/internal/model"
)

func (s *MemoryStore) CreateUsage(u *model.Usage) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.usages[u.ID] = u
	return nil
}

func (s *MemoryStore) GetUsage(id string) (*model.Usage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.usages[id]
	if !ok {
		return nil, &StoreError{Kind: KindNotFound, Op: "get usage", Err: ErrNotFound}
	}
	return u, nil
}

func (s *MemoryStore) ListUsages() []*model.Usage {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Usage, 0, len(s.usages))
	for _, u := range s.usages {
		list = append(list, u)
	}
	return list
}
