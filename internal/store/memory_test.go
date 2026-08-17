package store

import (
	"errors"
	"testing"
	"time"

	"coupon/internal/model"
)

func newTestStore() *MemoryStore { return NewMemoryStore() }

func testCoupon(name string) *model.Coupon {
	now := time.Now()
	return &model.Coupon{
		ID:         "c-" + name,
		Name:       name,
		Type:       model.CouponTypeFixed,
		Value:      500,
		TotalCount: 100,
		Status:     model.CouponActive,
		StartAt:    now.Add(-time.Hour),
		EndAt:      now.Add(24 * time.Hour),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func TestCouponCRUD(t *testing.T) {
	s := newTestStore()

	if err := s.CreateCoupon(testCoupon("满50减5")); err != nil {
		t.Fatalf("create: %v", err)
	}
	// 同名冲突
	if err := s.CreateCoupon(testCoupon("满50减5")); !errors.Is(err, ErrConflict) {
		t.Fatalf("expect conflict, got %v", err)
	}

	got, err := s.GetCoupon("c-满50减5")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Value != 500 {
		t.Fatalf("value = %d, want 500", got.Value)
	}

	if _, err := s.GetCoupon("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expect not found, got %v", err)
	}

	if n := len(s.ListCoupons()); n != 1 {
		t.Fatalf("list len = %d, want 1", n)
	}

	// 更新后值变化
	got.Value = 800
	if err := s.UpdateCoupon(got); err != nil {
		t.Fatalf("update: %v", err)
	}
	updated, _ := s.GetCoupon("c-满50减5")
	if updated.Value != 800 {
		t.Fatalf("after update value = %d, want 800", updated.Value)
	}

	if err := s.DeleteCoupon("c-满50减5"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if err := s.DeleteCoupon("c-满50减5"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expect not found on re-delete, got %v", err)
	}
}

func TestBatchCRUD(t *testing.T) {
	s := newTestStore()
	b := &model.Batch{ID: "b1", Name: "批次A", CouponID: "c1", TotalCount: 10}
	if err := s.CreateBatch(b); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := s.CreateBatch(&model.Batch{ID: "b2", Name: "批次A", CouponID: "c1", TotalCount: 5}); !errors.Is(err, ErrConflict) {
		t.Fatalf("expect conflict, got %v", err)
	}
	got, err := s.GetBatch("b1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.TotalCount != 10 {
		t.Fatalf("total = %d, want 10", got.TotalCount)
	}

	s.CreateBatch(&model.Batch{ID: "b3", Name: "批次B", CouponID: "c2", TotalCount: 3})
	byCoupon := s.ListBatchesByCoupon("c1")
	if len(byCoupon) != 1 {
		t.Fatalf("by coupon len = %d, want 1", len(byCoupon))
	}

	if err := s.DeleteBatch("b1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := s.GetBatch("b1"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expect not found, got %v", err)
	}
}

func TestUserCouponCRUD(t *testing.T) {
	s := newTestStore()
	uc := &model.UserCoupon{ID: "uc1", UserID: "u1", CouponID: "c1", Code: "ABC123", Status: model.UserCouponUnused}
	if err := s.CreateUserCoupon(uc); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := s.CreateUserCoupon(&model.UserCoupon{ID: "uc2", UserID: "u2", CouponID: "c1", Code: "ABC123"}); !errors.Is(err, ErrConflict) {
		t.Fatalf("expect conflict by code, got %v", err)
	}
	got, err := s.GetUserCouponByCode("ABC123")
	if err != nil {
		t.Fatalf("get by code: %v", err)
	}
	if got.UserID != "u1" {
		t.Fatalf("user = %s, want u1", got.UserID)
	}

	// 状态流转
	got.Status = model.UserCouponUsed
	if err := s.UpdateUserCoupon(got); err != nil {
		t.Fatalf("update: %v", err)
	}
	reloaded, _ := s.GetUserCoupon("uc1")
	if reloaded.Status != model.UserCouponUsed {
		t.Fatalf("status = %s, want used", reloaded.Status)
	}

	if err := s.DeleteUserCoupon("uc1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := s.GetUserCoupon("uc1"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expect not found, got %v", err)
	}
}

func TestUsageCRUD(t *testing.T) {
	s := newTestStore()
	u := &model.Usage{ID: "usage1", UserCouponID: "uc1", UserID: "u1", CouponID: "c1", Amount: 500}
	if err := s.CreateUsage(u); err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := s.GetUsage("usage1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Amount != 500 {
		t.Fatalf("amount = %d, want 500", got.Amount)
	}
	if _, err := s.GetUsage("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expect not found, got %v", err)
	}
	if n := len(s.ListUsages()); n != 1 {
		t.Fatalf("list len = %d, want 1", n)
	}
}
