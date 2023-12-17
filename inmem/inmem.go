package inmem

import (
	"context"
	"sync"

	"github.com/wizhi/zssn"
)

type SurvivorRepository struct {
	m    map[string]zssn.Survivor
	mu   sync.Mutex
	once sync.Once
}

func (r *SurvivorRepository) init() {
	r.once.Do(func() {
		r.m = make(map[string]zssn.Survivor)
	})
}

func (r *SurvivorRepository) Load(_ context.Context, id string) (*zssn.Survivor, error) {
	r.init()
	r.mu.Lock()
	defer r.mu.Unlock()

	s, ok := r.m[id]
	if !ok {
		return nil, zssn.ErrNotFound
	}

	return &s, nil
}

func (r *SurvivorRepository) Save(_ context.Context, ss ...*zssn.Survivor) error {
	r.init()
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, s := range ss {
		r.m[s.ID] = *s
	}

	return nil
}
