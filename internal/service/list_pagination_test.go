package service

import (
	"fmt"
	"testing"
	"time"

	"coupon/internal/config"
	"coupon/internal/model"
	"coupon/internal/store"
	"coupon/pkg/logger"
)

func paginationService() (*Service, *store.MemoryStore) {
	st := store.NewMemoryStore()
	svc := New(st, logger.New(), &config.Config{})
	return svc, st
}

func seedPaginationData(t *testing.T, svc *Service, st *store.MemoryStore) {
	t.Helper()
	now := time.Now()
	for i := 1; i <= 3; i++ {
		c := &model.Coupon{
			ID:         fmt.Sprintf("c%d", i),
			Name:       fmt.Sprintf("券%d", i),
			Type:       model.CouponTypeFixed,
			Value:      int64(100 * i),
			TotalCount: 10,
			Status:     model.CouponActive,
			StartAt:    now.Add(-time.Hour),
			EndAt:      now.Add(24 * time.Hour),
			CreatedAt:  now.Add(time.Duration(i) * time.Second),
			UpdatedAt:  now,
		}
		if err := st.CreateCoupon(c); err != nil {
			t.Fatal(err)
		}
		b := &model.Batch{
			ID:         fmt.Sprintf("b%d", i),
			Name:       fmt.Sprintf("批次%d", i),
			CouponID:   c.ID,
			TotalCount: 10,
			CreatedAt:  now.Add(time.Duration(i) * time.Second),
			UpdatedAt:  now,
		}
		if err := st.CreateBatch(b); err != nil {
			t.Fatal(err)
		}
		uc := &model.UserCoupon{
			ID:         fmt.Sprintf("u%d", i),
			UserID:     "u1",
			CouponID:   c.ID,
			BatchID:    b.ID,
			Code:       fmt.Sprintf("C%d", i),
			Status:     model.UserCouponUnused,
			ReceivedAt: now.Add(time.Duration(i) * time.Second),
		}
		if err := st.CreateUserCoupon(uc); err != nil {
			t.Fatal(err)
		}
		usage := &model.Usage{
			ID:           fmt.Sprintf("g%d", i),
			UserCouponID: uc.ID,
			UserID:       "u1",
			CouponID:     c.ID,
			Amount:       100,
			CreatedAt:    now.Add(time.Duration(i) * time.Second),
		}
		if err := st.CreateUsage(usage); err != nil {
			t.Fatal(err)
		}
	}
}

func TestListPagination(t *testing.T) {
	svc, st := paginationService()
	seedPaginationData(t, svc, st)

	t.Run("coupons", func(t *testing.T) {
		items, total, err := svc.ListCoupons(model.CouponFilter{}, 1, 2)
		if err != nil {
			t.Fatal(err)
		}
		if total != 3 || len(items) != 2 {
			t.Fatalf("coupons total=%d len=%d, want total=3 len=2", total, len(items))
		}
	})

	t.Run("batches", func(t *testing.T) {
		items, total, err := svc.ListBatches("", 1, 2)
		if err != nil {
			t.Fatal(err)
		}
		if total != 3 || len(items) != 2 {
			t.Fatalf("batches total=%d len=%d, want total=3 len=2", total, len(items))
		}
	})

	t.Run("user_coupons", func(t *testing.T) {
		items, total, err := svc.ListUserCoupons(model.UserCouponFilter{}, 1, 2)
		if err != nil {
			t.Fatal(err)
		}
		if total != 3 || len(items) != 2 {
			t.Fatalf("user_coupons total=%d len=%d, want total=3 len=2", total, len(items))
		}
	})

	t.Run("usages", func(t *testing.T) {
		items, total, err := svc.ListUsages(model.UsageFilter{}, 1, 2)
		if err != nil {
			t.Fatal(err)
		}
		if total != 3 || len(items) != 2 {
			t.Fatalf("usages total=%d len=%d, want total=3 len=2", total, len(items))
		}
	})

	// 翻页连续性：3 条数据按 size=2 分页，第一页 2 条、第二页 1 条、第三页空，
	// 且两页内容不重叠也不漏。
	t.Run("continuity", func(t *testing.T) {
		p1, total, err := svc.ListCoupons(model.CouponFilter{}, 1, 2)
		if err != nil {
			t.Fatal(err)
		}
		p2, _, err := svc.ListCoupons(model.CouponFilter{}, 2, 2)
		if err != nil {
			t.Fatal(err)
		}
		p3, _, err := svc.ListCoupons(model.CouponFilter{}, 3, 2)
		if err != nil {
			t.Fatal(err)
		}
		if total != 3 || len(p1) != 2 || len(p2) != 1 || len(p3) != 0 {
			t.Fatalf("continuity total=%d len(p1)=%d len(p2)=%d len(p3)=%d, want 3/2/1/0",
				total, len(p1), len(p2), len(p3))
		}
		seen := map[string]bool{}
		for _, c := range append(append(p1, p2...), p3...) {
			if seen[c.ID] {
				t.Fatalf("ID %s 在多页中重复出现", c.ID)
			}
			seen[c.ID] = true
		}
		if len(seen) != 3 {
			t.Fatalf("多页去重后仅 %d 条，want 3（存在漏页）", len(seen))
		}
	})
}
