// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"coupon/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// 优惠券
	CreateCoupon(c *model.Coupon) error
	GetCoupon(id string) (*model.Coupon, error)
	GetCouponByName(name string) (*model.Coupon, error)
	ListCoupons() []*model.Coupon
	UpdateCoupon(c *model.Coupon) error
	DeleteCoupon(id string) error

	// 批次
	CreateBatch(b *model.Batch) error
	GetBatch(id string) (*model.Batch, error)
	ListBatches() []*model.Batch
	ListBatchesByCoupon(couponID string) []*model.Batch
	UpdateBatch(b *model.Batch) error
	DeleteBatch(id string) error

	// 用户券
	IssueUserCoupon(userID, batchID string) (*model.UserCoupon, error)
	CreateUserCoupon(u *model.UserCoupon) error
	GetUserCoupon(id string) (*model.UserCoupon, error)
	GetUserCouponByCode(code string) (*model.UserCoupon, error)
	ListUserCoupons() []*model.UserCoupon
	UpdateUserCoupon(u *model.UserCoupon) error
	DeleteUserCoupon(id string) error

	// 使用记录
	CreateUsage(u *model.Usage) error
	GetUsage(id string) (*model.Usage, error)
	ListUsages() []*model.Usage
}
