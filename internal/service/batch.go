package service

import (
	"sort"
	"time"

	"coupon/internal/model"
	"coupon/internal/store"
	"coupon/pkg/idgen"
)

func (s *Service) CreateBatch(input model.Batch) (*model.Batch, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	// 跨实体校验：关联的券必须存在且处于 active 状态
	c, err := s.store.GetCoupon(input.CouponID)
	if err != nil {
		return nil, err
	}
	if c.Status != model.CouponActive {
		return nil, model.NewValidationError("coupon_id", "关联的券未启用")
	}
	now := time.Now()
	b := &model.Batch{
		ID:         idgen.Hex(),
		Name:       input.Name,
		CouponID:   input.CouponID,
		TotalCount: input.TotalCount,
		Status:     model.BatchCreated,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.store.CreateBatch(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) GetBatch(id string) (*model.Batch, error) {
	return s.store.GetBatch(id)
}

func (s *Service) ListBatches(couponID string, page, size int) ([]*model.Batch, int, error) {
	var all []*model.Batch
	if couponID != "" {
		all = s.store.ListBatchesByCoupon(couponID)
	} else {
		all = s.store.ListBatches()
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})
	total := len(all)
	start := page * size
	if start > total {
		return []*model.Batch{}, total, nil
	}
	end := start + size
	if end >= total {
		end = total
	}
	return all[start:end], total, nil
}

func (s *Service) DeleteBatch(id string) error {
	return s.store.DeleteBatch(id)
}

// startBatch 将批次置为发放中，供内部领取流程使用。
func (s *Service) ensureIssuing(b *model.Batch) error {
	if b.Status == model.BatchFinished {
		return store.ErrConflict
	}
	if b.Status == model.BatchCreated {
		b.Status = model.BatchIssuing
		b.UpdatedAt = time.Now()
		return s.store.UpdateBatch(b)
	}
	return nil
}
