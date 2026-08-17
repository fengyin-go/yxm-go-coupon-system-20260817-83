package handler

import (
	"net/http"
	"time"

	"coupon/internal/model"
	"coupon/pkg/httpx"
)

func (s *Server) registerCouponRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/coupons", s.createCoupon)
	mux.HandleFunc("GET /api/coupons", s.listCoupons)
	mux.HandleFunc("GET /api/coupons/{id}", s.getCoupon)
	mux.HandleFunc("PUT /api/coupons/{id}", s.updateCoupon)
	mux.HandleFunc("DELETE /api/coupons/{id}", s.deleteCoupon)
}

type createCouponRequest struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Value      int64  `json:"value"`
	MinSpend   int64  `json:"min_spend"`
	UserLimit  *int   `json:"user_limit"`
	TotalCount int    `json:"total_count"`
	Status     string `json:"status"`
	StartAt    string `json:"start_at"`
	EndAt      string `json:"end_at"`
}

func (s *Server) createCoupon(w http.ResponseWriter, r *http.Request) {
	var req createCouponRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	startAt, err := parseTime(req.StartAt)
	if err != nil {
		httpx.BadRequest(w, "start_at 时间格式非法")
		return
	}
	endAt, err := parseTime(req.EndAt)
	if err != nil {
		httpx.BadRequest(w, "end_at 时间格式非法")
		return
	}
	c, err := s.svc.CreateCoupon(model.Coupon{
		Name:       req.Name,
		Type:       req.Type,
		Value:      req.Value,
		MinSpend:   req.MinSpend,
		UserLimit:  req.UserLimit,
		TotalCount: req.TotalCount,
		Status:     req.Status,
		StartAt:    startAt,
		EndAt:      endAt,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, c)
}

func (s *Server) listCoupons(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.CouponFilter{
		Type:    r.URL.Query().Get("type"),
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListCoupons(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getCoupon(w http.ResponseWriter, r *http.Request) {
	c, err := s.svc.GetCoupon(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) updateCoupon(w http.ResponseWriter, r *http.Request) {
	var req createCouponRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	input := model.Coupon{
		Name:      req.Name,
		Type:      req.Type,
		Value:     req.Value,
		MinSpend:  req.MinSpend,
		UserLimit: req.UserLimit,
		Status:    req.Status,
	}
	if req.StartAt != "" {
		t, err := parseTime(req.StartAt)
		if err != nil {
			httpx.BadRequest(w, "start_at 时间格式非法")
			return
		}
		input.StartAt = t
	}
	if req.EndAt != "" {
		t, err := parseTime(req.EndAt)
		if err != nil {
			httpx.BadRequest(w, "end_at 时间格式非法")
			return
		}
		input.EndAt = t
	}
	c, err := s.svc.UpdateCoupon(r.PathValue("id"), input)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) deleteCoupon(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteCoupon(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

// parseTime 解析 RFC3339 或 "2006-01-02 15:04:05" 格式的时间。
func parseTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	return time.Parse("2006-01-02 15:04:05", s)
}
