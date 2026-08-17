package store

import (
	"time"

	"coupon/internal/model"
	"coupon/pkg/idgen"
)

// IssueUserCoupon 在单把写锁内完成领取流程的检查与计数更新，避免并发超发。
func (s *MemoryStore) IssueUserCoupon(userID, batchID string) (*model.UserCoupon, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if userID == "" {
		return nil, model.NewValidationError("user_id", "用户不能为空")
	}
	b, ok := s.batches[batchID]
	if !ok {
		return nil, ErrNotFound
	}
	c, ok := s.coupons[b.CouponID]
	if !ok {
		return nil, ErrNotFound
	}
	if c.Status != model.CouponActive {
		return nil, model.NewValidationError("coupon", "券未启用")
	}
	if !c.InWindow(time.Now()) {
		return nil, model.NewValidationError("coupon", "券不在有效期内")
	}
	if !b.Reserve() {
		return nil, ErrConflict
	}

	now := time.Now()
	uc := &model.UserCoupon{
		ID:         idgen.Hex(),
		UserID:     userID,
		CouponID:   c.ID,
		BatchID:    b.ID,
		Code:       idgen.Short(),
		Status:     model.UserCouponUnused,
		ReceivedAt: now,
		ExpiredAt:  c.EndAt,
	}
	s.userCoupons[uc.ID] = uc
	c.IssuedCount++
	return uc, nil
}
