package service

import (
	"sort"

	"coupon/internal/model"
)

func (s *Service) GetUsage(id string) (*model.Usage, error) {
	return s.store.GetUsage(id)
}

func (s *Service) ListUsages(filter model.UsageFilter, page, size int) ([]*model.Usage, int, error) {
	all := s.store.ListUsages()
	matched := make([]*model.Usage, 0, len(all))
	for _, u := range all {
		if filter.Match(u) {
			matched = append(matched, u)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Usage{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}
