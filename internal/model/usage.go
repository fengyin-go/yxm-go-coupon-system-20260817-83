package model

import (
	"strings"
	"time"
)

// Usage 券核销使用记录，一条用户券被使用即生成一条记录。
// Amount 表示本次实际抵扣金额，单位分。
type Usage struct {
	ID           string    `json:"id"`
	UserCouponID string    `json:"user_coupon_id"`
	UserID       string    `json:"user_id"`
	CouponID     string    `json:"coupon_id"`
	OrderID      string    `json:"order_id"`
	Amount       int64     `json:"amount"`
	CreatedAt    time.Time `json:"created_at"`
}

func (u *Usage) Validate() error {
	u.UserCouponID = strings.TrimSpace(u.UserCouponID)
	u.UserID = strings.TrimSpace(u.UserID)
	if u.UserCouponID == "" {
		return NewValidationError("user_coupon_id", "用户券不能为空")
	}
	if u.UserID == "" {
		return NewValidationError("user_id", "用户不能为空")
	}
	if u.Amount <= 0 {
		return NewValidationError("amount", "抵扣金额必须大于 0")
	}
	return nil
}

type UsageFilter struct {
	UserID   string
	CouponID string
}

func (f UsageFilter) Match(u *Usage) bool {
	if f.UserID != "" && u.UserID != f.UserID {
		return false
	}
	if f.CouponID != "" && u.CouponID != f.CouponID {
		return false
	}
	return true
}
