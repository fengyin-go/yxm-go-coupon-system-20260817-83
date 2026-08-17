package service

import (
	"sort"
	"time"

	"coupon/internal/model"
	"coupon/internal/store"
	"coupon/pkg/idgen"
)

// Issue 从指定批次向用户发放一张券，完成批次/券的发放计数更新。
func (s *Service) Issue(userID, batchID string) (*model.UserCoupon, error) {
	if userID == "" {
		return nil, model.NewValidationError("user_id", "用户不能为空")
	}
	b, err := s.store.GetBatch(batchID)
	if err != nil {
		return nil, err
	}
	c, err := s.store.GetCoupon(b.CouponID)
	if err != nil {
		return nil, err
	}
	if c.Status != model.CouponActive {
		return nil, model.NewValidationError("coupon", "券未启用")
	}
	if !c.InWindow(time.Now()) {
		return nil, model.NewValidationError("coupon", "券不在有效期内")
	}
	if b.Remaining() <= 0 {
		return nil, store.ErrConflict
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
	if err := s.store.CreateUserCoupon(uc); err != nil {
		return nil, err
	}

	b.IssuedCount++
	if b.IssuedCount >= b.TotalCount {
		b.Status = model.BatchFinished
	} else if b.Status == model.BatchCreated {
		b.Status = model.BatchIssuing
	}
	b.UpdatedAt = now
	if err := s.store.UpdateBatch(b); err != nil {
		return nil, err
	}

	c.IssuedCount++
	if err := s.store.UpdateCoupon(c); err != nil {
		return nil, err
	}
	return uc, nil
}

// Use 核销一张用户券，完成状态机流转并生成使用记录。
func (s *Service) Use(userCouponID, orderID string) (*model.Usage, error) {
	uc, err := s.store.GetUserCoupon(userCouponID)
	if err != nil {
		return nil, err
	}
	if uc.Status != model.UserCouponUnused {
		return nil, store.ErrConflict
	}
	c, err := s.store.GetCoupon(uc.CouponID)
	if err != nil {
		return nil, err
	}
	if c.Status != model.CouponActive {
		return nil, model.NewValidationError("coupon", "券未启用")
	}
	if !c.InWindow(time.Now()) {
		return nil, model.NewValidationError("coupon", "券已过期")
	}

	now := time.Now()
	// 状态机：unused -> used
	if !model.CanTransitionUserCoupon(uc.Status, model.UserCouponUsed) {
		return nil, store.ErrConflict
	}
	uc.Status = model.UserCouponUsed
	uc.UsedAt = &now
	if err := s.store.UpdateUserCoupon(uc); err != nil {
		return nil, err
	}

	usage := &model.Usage{
		ID:           idgen.Hex(),
		UserCouponID: uc.ID,
		UserID:       uc.UserID,
		CouponID:     uc.CouponID,
		OrderID:      orderID,
		Amount:       c.Value,
		CreatedAt:    now,
	}
	if err := s.store.CreateUsage(usage); err != nil {
		return nil, err
	}

	c.UsedCount++
	if err := s.store.UpdateCoupon(c); err != nil {
		return nil, err
	}
	return usage, nil
}

func (s *Service) GetUserCoupon(id string) (*model.UserCoupon, error) {
	return s.store.GetUserCoupon(id)
}

func (s *Service) ListUserCoupons(filter model.UserCouponFilter, page, size int) ([]*model.UserCoupon, int, error) {
	all := s.store.ListUserCoupons()
	matched := make([]*model.UserCoupon, 0, len(all))
	for _, u := range all {
		if filter.Match(u) {
			matched = append(matched, u)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].ReceivedAt.After(matched[j].ReceivedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.UserCoupon{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}
