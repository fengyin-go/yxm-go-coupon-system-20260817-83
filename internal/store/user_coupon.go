package store

import (
	"coupon/internal/model"
)

func (s *MemoryStore) CreateUserCoupon(u *model.UserCoupon) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.userCoupons {
		if exist.Code == u.Code {
			return &StoreError{Kind: KindConflict, Op: "create user coupon", Err: ErrConflict}
		}
	}
	s.userCoupons[u.ID] = u
	return nil
}

func (s *MemoryStore) GetUserCoupon(id string) (*model.UserCoupon, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.userCoupons[id]
	if !ok {
		return nil, &StoreError{Kind: KindNotFound, Op: "get user coupon", Err: ErrNotFound}
	}
	return u, nil
}

func (s *MemoryStore) GetUserCouponByCode(code string) (*model.UserCoupon, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.userCoupons {
		if u.Code == code {
			return u, nil
		}
	}
	return nil, &StoreError{Kind: KindNotFound, Op: "get user coupon by code", Err: ErrNotFound}
}

func (s *MemoryStore) ListUserCoupons() []*model.UserCoupon {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.UserCoupon, 0, len(s.userCoupons))
	for _, u := range s.userCoupons {
		list = append(list, u)
	}
	return list
}

func (s *MemoryStore) UpdateUserCoupon(u *model.UserCoupon) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.userCoupons[u.ID]; !ok {
		return &StoreError{Kind: KindNotFound, Op: "update user coupon", Err: ErrNotFound}
	}
	s.userCoupons[u.ID] = u
	return nil
}

func (s *MemoryStore) DeleteUserCoupon(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.userCoupons[id]; !ok {
		return &StoreError{Kind: KindNotFound, Op: "delete user coupon", Err: ErrNotFound}
	}
	delete(s.userCoupons, id)
	return nil
}
