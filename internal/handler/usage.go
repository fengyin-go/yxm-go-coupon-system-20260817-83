package handler

import (
	"net/http"

	"coupon/internal/model"
	"coupon/pkg/httpx"
)

func (s *Server) registerUsageRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/usages", s.listUsages)
	mux.HandleFunc("GET /api/usages/{id}", s.getUsage)
}

func (s *Server) listUsages(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.UsageFilter{
		UserID:   r.URL.Query().Get("user_id"),
		CouponID: r.URL.Query().Get("coupon_id"),
	}
	items, total, err := s.svc.ListUsages(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getUsage(w http.ResponseWriter, r *http.Request) {
	u, err := s.svc.GetUsage(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, u)
}
