package store

import (
	"coupon/internal/model"
)

func (s *MemoryStore) CreateCoupon(c *model.Coupon) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.coupons {
		if exist.Name == c.Name {
			return &StoreError{Kind: KindConflict, Op: "create coupon", Err: ErrConflict}
		}
	}
	s.coupons[c.ID] = c
	return nil
}

func (s *MemoryStore) GetCoupon(id string) (*model.Coupon, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.coupons[id]
	if !ok {
		return nil, &StoreError{Kind: KindNotFound, Op: "get coupon", Err: ErrNotFound}
	}
	return c, nil
}

func (s *MemoryStore) GetCouponByName(name string) (*model.Coupon, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.coupons {
		if c.Name == name {
			return c, nil
		}
	}
	return nil, &StoreError{Kind: KindNotFound, Op: "get coupon by name", Err: ErrNotFound}
}

func (s *MemoryStore) ListCoupons() []*model.Coupon {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Coupon, 0, len(s.coupons))
	for _, c := range s.coupons {
		list = append(list, c)
	}
	return list
}

func (s *MemoryStore) UpdateCoupon(c *model.Coupon) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.coupons[c.ID]; !ok {
		return &StoreError{Kind: KindNotFound, Op: "update coupon", Err: ErrNotFound}
	}
	for _, exist := range s.coupons {
		if exist.ID != c.ID && exist.Name == c.Name {
			return &StoreError{Kind: KindConflict, Op: "update coupon", Err: ErrConflict}
		}
	}
	s.coupons[c.ID] = c
	return nil
}

func (s *MemoryStore) DeleteCoupon(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.coupons[id]; !ok {
		return &StoreError{Kind: KindNotFound, Op: "delete coupon", Err: ErrNotFound}
	}
	delete(s.coupons, id)
	return nil
}
