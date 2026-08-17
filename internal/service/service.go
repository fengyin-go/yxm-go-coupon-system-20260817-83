// Package service 实现业务逻辑层。
package service

import (
	"coupon/internal/config"
	"coupon/internal/store"
	"coupon/pkg/logger"
)

type Service struct {
	store store.Store
	log   *logger.Logger
	cfg   *config.Config
}

func New(st store.Store, log *logger.Logger, cfg *config.Config) *Service {
	return &Service{store: st, log: log, cfg: cfg}
}
