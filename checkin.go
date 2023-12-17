package zssn

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CheckInCommand struct {
	SurvivorID          string
	Latitude, Longitude float64
}

type CheckInHandler struct {
	Survivors SurvivorRepository
}

func (h *CheckInHandler) Handle(ctx context.Context, cmd CheckInCommand) error {
	s, err := h.Survivors.Load(ctx, cmd.SurvivorID)
	if err != nil {
		return err
	}

	s.CheckIn(cmd.Latitude, cmd.Longitude)

	return h.Survivors.Save(ctx, s)
}

func (h *CheckInHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var cmd CheckInCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		log.Printf("failed to decode request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cmd.SurvivorID = chi.URLParam(r, "survivorID")

	err := h.Handle(r.Context(), cmd)
	switch {
	case err == nil:
		w.WriteHeader(http.StatusNoContent)
	case errors.Is(err, ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
