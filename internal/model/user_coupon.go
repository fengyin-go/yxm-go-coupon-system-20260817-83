package model

import (
	"strings"
	"time"
)

const (
	UserCouponUnused  = "unused"
	UserCouponUsed    = "used"
	UserCouponExpired = "expired"
)

// userCouponTransitions 定义用户券状态机：unused 可核销为 used，也可过期为 expired。
var userCouponTransitions = map[string]map[string]bool{
	UserCouponUnused: {UserCouponUsed: true, UserCouponExpired: true},
}

// CanTransition 判断用户券状态流转是否合法。
func CanTransitionUserCoupon(from, to string) bool {
	if m, ok := userCouponTransitions[from]; ok {
		return m[to]
	}
	return false
}

// UserCoupon 用户持有的券实例，由批次发放给具体用户。
type UserCoupon struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	CouponID   string     `json:"coupon_id"`
	BatchID    string     `json:"batch_id"`
	Code       string     `json:"code"`
	Status     string     `json:"status"`
	ReceivedAt time.Time  `json:"received_at"`
	UsedAt     *time.Time `json:"used_at,omitempty"`
	ExpiredAt  time.Time  `json:"expired_at"`
}

func (u *UserCoupon) Validate() error {
	u.UserID = strings.TrimSpace(u.UserID)
	if u.UserID == "" {
		return NewValidationError("user_id", "用户不能为空")
	}
	if u.CouponID == "" {
		return NewValidationError("coupon_id", "关联券不能为空")
	}
	if u.Status == "" {
		u.Status = UserCouponUnused
	}
	if u.Status != UserCouponUnused && u.Status != UserCouponUsed && u.Status != UserCouponExpired {
		return NewValidationError("status", "用户券状态不合法")
	}
	return nil
}

type UserCouponFilter struct {
	UserID string
	Status string
}

func (f UserCouponFilter) Match(u *UserCoupon) bool {
	if f.UserID != "" && u.UserID != f.UserID {
		return false
	}
	if f.Status != "" && u.Status != f.Status {
		return false
	}
	return true
}
