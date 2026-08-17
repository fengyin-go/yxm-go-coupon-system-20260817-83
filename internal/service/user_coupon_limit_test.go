package service

import (
	"testing"
	"time"

	"coupon/internal/config"
	"coupon/internal/model"
	"coupon/internal/store"
	"coupon/pkg/logger"
)

func newTestService() (*Service, *store.MemoryStore) {
	st := store.NewMemoryStore()
	svc := New(st, logger.New(), &config.Config{})
	return svc, st
}

func mustCreateCouponForLimit(t *testing.T, st *store.MemoryStore, userLimit *int) *model.Coupon {
	t.Helper()
	now := time.Now()
	c := &model.Coupon{
		ID:         "c-limit",
		Name:       "限领券",
		Type:       model.CouponTypeFixed,
		Value:      500,
		MinSpend:   0,
		UserLimit:  userLimit,
		TotalCount: 10,
		Status:     model.CouponActive,
		StartAt:    now.Add(-time.Hour),
		EndAt:      now.Add(24 * time.Hour),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := st.CreateCoupon(c); err != nil {
		t.Fatalf("create coupon: %v", err)
	}
	return c
}

func TestIssueCouponWithoutUserLimitDoesNotPanic(t *testing.T) {
	svc, st := newTestService()
	mustCreateCouponForLimit(t, st, nil)
	if err := st.CreateBatch(&model.Batch{ID: "b-no-limit", Name: "无限领批次", CouponID: "c-limit", TotalCount: 10}); err != nil {
		t.Fatalf("create batch: %v", err)
	}

	uc, err := svc.Issue("u1", "b-no-limit")
	if err != nil {
		t.Fatalf("issue without user_limit should succeed, got %v", err)
	}
	if uc == nil || uc.Status != model.UserCouponUnused {
		t.Fatalf("unexpected user coupon: %#v", uc)
	}
}

func TestIssueEnforcesUserLimit(t *testing.T) {
	svc, st := newTestService()
	limit := 1
	mustCreateCouponForLimit(t, st, &limit)
	if err := st.CreateBatch(&model.Batch{ID: "b-limit", Name: "限领批次", CouponID: "c-limit", TotalCount: 10}); err != nil {
		t.Fatalf("create batch: %v", err)
	}

	if _, err := svc.Issue("u1", "b-limit"); err != nil {
		t.Fatalf("first issue should succeed, got %v", err)
	}
	if _, err := svc.Issue("u1", "b-limit"); err != store.ErrConflict {
		t.Fatalf("second issue should be conflict, got %v", err)
	}
}
