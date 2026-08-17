package store

import (
	"sync"

	"coupon/internal/model"
)

// MemoryStore 基于内存的 Store 实现，使用读写锁保证并发安全。
type MemoryStore struct {
	mu          sync.RWMutex
	coupons     map[string]*model.Coupon
	batches     map[string]*model.Batch
	userCoupons map[string]*model.UserCoupon
	usages      map[string]*model.Usage
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		coupons:     make(map[string]*model.Coupon),
		batches:     make(map[string]*model.Batch),
		userCoupons: make(map[string]*model.UserCoupon),
		usages:      make(map[string]*model.Usage),
	}
}

var _ Store = (*MemoryStore)(nil)
