package handler

import (
	"net/http"

	"coupon/internal/model"
	"coupon/pkg/httpx"
)

func (s *Server) registerUserCouponRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/user-coupons/issue", s.issueUserCoupon)
	mux.HandleFunc("POST /api/user-coupons/{id}/use", s.useUserCoupon)
	mux.HandleFunc("GET /api/user-coupons", s.listUserCoupons)
	mux.HandleFunc("GET /api/user-coupons/{id}", s.getUserCoupon)
}

type issueUserCouponRequest struct {
	UserID  string `json:"user_id"`
	BatchID string `json:"batch_id"`
}

func (s *Server) issueUserCoupon(w http.ResponseWriter, r *http.Request) {
	var req issueUserCouponRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	uc, err := s.svc.Issue(req.UserID, req.BatchID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, uc)
}

type useUserCouponRequest struct {
	OrderID string `json:"order_id"`
}

func (s *Server) useUserCoupon(w http.ResponseWriter, r *http.Request) {
	var req useUserCouponRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	usage, err := s.svc.Use(r.PathValue("id"), req.OrderID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, usage)
}

func (s *Server) listUserCoupons(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.UserCouponFilter{
		UserID: r.URL.Query().Get("user_id"),
		Status: r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListUserCoupons(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getUserCoupon(w http.ResponseWriter, r *http.Request) {
	uc, err := s.svc.GetUserCoupon(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, uc)
}
