package zssn

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type StatusQuery struct {
	SurvivorID string
}

type StatusResult struct {
	Name      string
	Age       int
	Gender    Gender
	Location  Location
	Infected  bool
	Inventory map[string]int
}

type StatusHandler struct {
	Survivors SurvivorRepository
}

func (h *StatusHandler) Handle(ctx context.Context, q StatusQuery) (StatusResult, error) {
	s, err := h.Survivors.Load(ctx, q.SurvivorID)
	if err != nil {
		return StatusResult{}, err
	}

	r := make(map[string]int, len(s.Inventory))
	for item, q := range s.Inventory {
		r[item.Kind] = q
	}

	return StatusResult{
		Name:      s.Name,
		Age:       s.Age,
		Gender:    s.Gender,
		Location:  s.Location,
		Infected:  s.Infected(),
		Inventory: r,
	}, nil
}

func (h *StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "survivorID")

	res, err := h.Handle(r.Context(), StatusQuery{SurvivorID: id})
	switch {
	case err == nil:
		if err := json.NewEncoder(w).Encode(res); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	case errors.Is(err, ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
