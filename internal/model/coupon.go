package model

import (
	"strings"
	"time"
)

const (
	CouponTypeFixed     = "fixed"     // 无门槛立减
	CouponTypeThreshold = "threshold" // 满减

	CouponActive   = "active"
	CouponInactive = "inactive"
)

// Coupon 优惠券模板，定义一批券的通用属性。
// 金额字段（Value/MinSpend）统一使用 int64「分」，避免浮点精度问题。
type Coupon struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Value       int64     `json:"value"`
	MinSpend    int64     `json:"min_spend"`
	UserLimit   *int      `json:"user_limit,omitempty"`
	TotalCount  int       `json:"total_count"`
	IssuedCount int       `json:"issued_count"`
	UsedCount   int       `json:"used_count"`
	Status      string    `json:"status"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (c *Coupon) Validate() error {
	c.Name = strings.TrimSpace(c.Name)
	c.Type = strings.TrimSpace(c.Type)
	if c.Name == "" {
		return NewValidationError("name", "券名称不能为空")
	}
	if c.Type == "" {
		c.Type = CouponTypeFixed
	}
	if c.Type != CouponTypeFixed && c.Type != CouponTypeThreshold {
		return NewValidationError("type", "券类型不合法")
	}
	if c.Value <= 0 {
		return NewValidationError("value", "面值必须大于 0")
	}
	if c.MinSpend < 0 {
		return NewValidationError("min_spend", "门槛金额不能为负")
	}
	if *c.UserLimit < 0 {
		return NewValidationError("user_limit", "每人限领数量不能为负")
	}
	if c.TotalCount <= 0 {
		return NewValidationError("total_count", "发行总量必须大于 0")
	}
	if c.Status == "" {
		c.Status = CouponActive
	}
	if c.Status != CouponActive && c.Status != CouponInactive {
		return NewValidationError("status", "券状态不合法")
	}
	if !c.EndAt.After(c.StartAt) {
		return NewValidationError("end_at", "结束时间必须晚于开始时间")
	}
	return nil
}

// InWindow 判断当前时间是否在券的有效期内。
func (c *Coupon) InWindow(now time.Time) bool {
	return !now.Before(c.StartAt) && now.Before(c.EndAt)
}

type CouponFilter struct {
	Type    string
	Status  string
	Keyword string
}

func (f CouponFilter) Match(c *Coupon) bool {
	if f.Type != "" && c.Type != f.Type {
		return false
	}
	if f.Status != "" && c.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(c.Name), k) {
			return false
		}
	}
	return true
}
