package model

import (
	"strings"
	"time"
)

const (
	BatchCreated  = "created"
	BatchIssuing  = "issuing"
	BatchFinished = "finished"
)

// Batch 券的发放批次，关联一张券模板，记录本次发放的数量。
type Batch struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CouponID    string    `json:"coupon_id"`
	TotalCount  int       `json:"total_count"`
	IssuedCount int       `json:"issued_count"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (b *Batch) Validate() error {
	b.Name = strings.TrimSpace(b.Name)
	if b.Name == "" {
		return NewValidationError("name", "批次名称不能为空")
	}
	if b.CouponID == "" {
		return NewValidationError("coupon_id", "关联券不能为空")
	}
	if b.TotalCount <= 0 {
		return NewValidationError("total_count", "发放数量必须大于 0")
	}
	if b.Status == "" {
		b.Status = BatchCreated
	}
	if b.Status != BatchCreated && b.Status != BatchIssuing && b.Status != BatchFinished {
		return NewValidationError("status", "批次状态不合法")
	}
	return nil
}

// Remaining 返回批次剩余可发放数量。
func (b *Batch) Remaining() int {
	if b.IssuedCount >= b.TotalCount {
		return 0
	}
	return b.TotalCount - b.IssuedCount
}
