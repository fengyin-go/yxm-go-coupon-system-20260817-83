package handler

import (
	"net/http"

	"coupon/internal/model"
	"coupon/pkg/httpx"
)

func (s *Server) registerBatchRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/batches", s.createBatch)
	mux.HandleFunc("GET /api/batches", s.listBatches)
	mux.HandleFunc("GET /api/batches/{id}", s.getBatch)
	mux.HandleFunc("DELETE /api/batches/{id}", s.deleteBatch)
}

type createBatchRequest struct {
	Name       string `json:"name"`
	CouponID   string `json:"coupon_id"`
	TotalCount int    `json:"total_count"`
}

func (s *Server) createBatch(w http.ResponseWriter, r *http.Request) {
	var req createBatchRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	b, err := s.svc.CreateBatch(model.Batch{
		Name:       req.Name,
		CouponID:   req.CouponID,
		TotalCount: req.TotalCount,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, b)
}

func (s *Server) listBatches(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	items, total, err := s.svc.ListBatches(r.URL.Query().Get("coupon_id"), pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getBatch(w http.ResponseWriter, r *http.Request) {
	b, err := s.svc.GetBatch(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}

func (s *Server) deleteBatch(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteBatch(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
