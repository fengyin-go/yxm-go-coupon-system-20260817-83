// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"
	"fmt"

	"coupon/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

type Kind int

const (
	KindNotFound Kind = iota
	KindConflict
)

// StoreError 包装 store 层错误，携带操作名与错误类别；
// Is 方法使 errors.Is(err, ErrNotFound/ErrConflict) 能正确匹配包装后的错误。
type StoreError struct {
	Kind Kind
	Op   string
	Err  error
}

func (e *StoreError) Error() string {
	return fmt.Sprintf("%s: %v", e.Op, e.Err)
}

func (e *StoreError) Is(target error) bool {
	switch target {
	case ErrNotFound:
		return e.Kind == KindNotFound
	case ErrConflict:
		return e.Kind == KindConflict
	default:
		return false
	}
}

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
