package service

import (
	"sort"
	"time"

	"coupon/internal/model"
	"coupon/pkg/idgen"
)

func (s *Service) CreateCoupon(input model.Coupon) (*model.Coupon, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	c := &model.Coupon{
		ID:         idgen.Hex(),
		Name:       input.Name,
		Type:       input.Type,
		Value:      input.Value,
		MinSpend:   input.MinSpend,
		TotalCount: input.TotalCount,
		Status:     input.Status,
		StartAt:    input.StartAt,
		EndAt:      input.EndAt,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.store.CreateCoupon(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) GetCoupon(id string) (*model.Coupon, error) {
	return s.store.GetCoupon(id)
}

func (s *Service) ListCoupons(filter model.CouponFilter, page, size int) ([]*model.Coupon, int, error) {
	all := s.store.ListCoupons()
	matched := make([]*model.Coupon, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := page * size
	if start > total {
		return []*model.Coupon{}, total, nil
	}
	end := start + size
	if end >= total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateCoupon(id string, input model.Coupon) (*model.Coupon, error) {
	existing, err := s.store.GetCoupon(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		existing.Name = input.Name
	}
	if input.Type != "" {
		existing.Type = input.Type
	}
	if input.Value > 0 {
		existing.Value = input.Value
	}
	if input.MinSpend >= 0 {
		existing.MinSpend = input.MinSpend
	}
	if input.Status != "" {
		existing.Status = input.Status
	}
	if !input.StartAt.IsZero() {
		existing.StartAt = input.StartAt
	}
	if !input.EndAt.IsZero() {
		existing.EndAt = input.EndAt
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateCoupon(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteCoupon(id string) error {
	return s.store.DeleteCoupon(id)
}
