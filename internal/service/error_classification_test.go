package service

import (
	"errors"
	"testing"
	"time"

	"coupon/internal/config"
	"coupon/internal/model"
	"coupon/internal/store"
	"coupon/pkg/logger"
)

func errorService() (*Service, *store.MemoryStore) {
	st := store.NewMemoryStore()
	svc := New(st, logger.New(), &config.Config{})
	return svc, st
}

func TestNotFoundErrorsAreClassified(t *testing.T) {
	svc, _ := errorService()

	if _, err := svc.GetCoupon("missing"); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("coupon missing error = %v, want ErrNotFound", err)
	}
	if _, err := svc.GetBatch("missing"); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("batch missing error = %v, want ErrNotFound", err)
	}
	if _, err := svc.GetUserCoupon("missing"); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("user coupon missing error = %v, want ErrNotFound", err)
	}
	if _, err := svc.GetUsage("missing"); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("usage missing error = %v, want ErrNotFound", err)
	}
}

func TestConflictErrorsAreClassified(t *testing.T) {
	svc, st := errorService()
	now := time.Now()

	c := &model.Coupon{
		ID:         "c1",
		Name:       "重复券",
		Type:       model.CouponTypeFixed,
		Value:      500,
		TotalCount: 10,
		Status:     model.CouponActive,
		StartAt:    now.Add(-time.Hour),
		EndAt:      now.Add(24 * time.Hour),
	}
	if err := st.CreateCoupon(c); err != nil {
		t.Fatalf("create coupon: %v", err)
	}

	dup := model.Coupon{
		Name:       c.Name,
		Type:       c.Type,
		Value:      c.Value,
		MinSpend:   c.MinSpend,
		TotalCount: c.TotalCount,
		Status:     c.Status,
		StartAt:    c.StartAt,
		EndAt:      c.EndAt,
	}
	if _, err := svc.CreateCoupon(dup); !errors.Is(err, store.ErrConflict) {
		t.Fatalf("duplicate coupon error = %v, want ErrConflict", err)
	}
}
