package store

import (
	"coupon/internal/model"
)

func (s *MemoryStore) CreateUserCoupon(u *model.UserCoupon) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.userCoupons {
		if exist.Code == u.Code {
			return ErrConflict
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
		return nil, ErrNotFound
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
	return nil, ErrNotFound
}

func (s *MemoryStore) CountUserCouponsByCouponAndUser(couponID, userID string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n := 0
	for _, u := range s.userCoupons {
		if u.UserID == userID {
			n++
		}
	}
	return n, nil
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
		return ErrNotFound
	}
	s.userCoupons[u.ID] = u
	return nil
}

func (s *MemoryStore) DeleteUserCoupon(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.userCoupons[id]; !ok {
		return ErrNotFound
	}
	delete(s.userCoupons, id)
	return nil
}
