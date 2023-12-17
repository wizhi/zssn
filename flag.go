package zssn

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type FlagCommand struct {
	SurvivorID string
}

type FlagHandler struct {
	Survivors SurvivorRepository
}

func (h *FlagHandler) Handle(ctx context.Context, cmd FlagCommand) error {
	s, err := h.Survivors.Load(ctx, cmd.SurvivorID)
	if err != nil {
		return err
	}

	s.Flag()

	return h.Survivors.Save(ctx, s)
}

func (h *FlagHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.Handle(r.Context(), FlagCommand{
		SurvivorID: chi.URLParam(r, "survivorID"),
	})
	switch {
	case err == nil:
		w.WriteHeader(http.StatusNoContent)
	case errors.Is(err, ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
