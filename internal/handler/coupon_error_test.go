package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"coupon/internal/config"
	"coupon/internal/service"
	"coupon/internal/store"
	"coupon/pkg/logger"
)

func newTestServer() *Server {
	st := store.NewMemoryStore()
	svc := service.New(st, logger.New(), &config.Config{})
	return NewServer(svc, logger.New(), &config.Config{})
}

const validCouponBody = `{"name":"测试券","type":"fixed","value":500,"total_count":10,"start_at":"2026-01-01T00:00:00Z","end_at":"2026-12-31T00:00:00Z"}`

// 请求不存在的券应返回 404，而非被误判为冲突（409）或内部错误（500）。
func TestGetMissingCouponReturns404(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/api/coupons/missing", nil)
	rr := httptest.NewRecorder()

	srv.Routes().ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("GET missing coupon status = %d, want 404; body=%s", rr.Code, rr.Body.String())
	}
}

// 创建重名券应返回 409，而非被误判为不存在（404）或内部错误（500）。
func TestCreateDuplicateCouponReturns409(t *testing.T) {
	srv := newTestServer()
	body := []byte(validCouponBody)

	req1 := httptest.NewRequest(http.MethodPost, "/api/coupons", bytes.NewReader(body))
	rr1 := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rr1, req1)
	if rr1.Code != http.StatusCreated {
		t.Fatalf("first create status = %d, want 201; body=%s", rr1.Code, rr1.Body.String())
	}

	req2 := httptest.NewRequest(http.MethodPost, "/api/coupons", bytes.NewReader(body))
	rr2 := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusConflict {
		t.Fatalf("duplicate create status = %d, want 409; body=%s", rr2.Code, rr2.Body.String())
	}
}
