package store

import (
	"coupon/internal/model"
)

func (s *MemoryStore) CreateBatch(b *model.Batch) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.batches {
		if exist.Name == b.Name {
			return &StoreError{Kind: KindConflict, Op: "create batch", Err: ErrConflict}
		}
	}
	s.batches[b.ID] = b
	return nil
}

func (s *MemoryStore) GetBatch(id string) (*model.Batch, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.batches[id]
	if !ok {
		return nil, &StoreError{Kind: KindNotFound, Op: "get batch", Err: ErrNotFound}
	}
	return b, nil
}

func (s *MemoryStore) ListBatches() []*model.Batch {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Batch, 0, len(s.batches))
	for _, b := range s.batches {
		list = append(list, b)
	}
	return list
}

func (s *MemoryStore) ListBatchesByCoupon(couponID string) []*model.Batch {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Batch, 0)
	for _, b := range s.batches {
		if b.CouponID == couponID {
			list = append(list, b)
		}
	}
	return list
}

func (s *MemoryStore) UpdateBatch(b *model.Batch) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.batches[b.ID]; !ok {
		return &StoreError{Kind: KindNotFound, Op: "update batch", Err: ErrNotFound}
	}
	for _, exist := range s.batches {
		if exist.ID != b.ID && exist.Name == b.Name {
			return &StoreError{Kind: KindConflict, Op: "update batch", Err: ErrConflict}
		}
	}
	s.batches[b.ID] = b
	return nil
}

func (s *MemoryStore) DeleteBatch(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.batches[id]; !ok {
		return &StoreError{Kind: KindNotFound, Op: "delete batch", Err: ErrNotFound}
	}
	delete(s.batches, id)
	return nil
}
