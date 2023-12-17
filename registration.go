package zssn

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rs/xid"
)

type RegistrationCommand struct {
	Name      string
	Age       int
	Gender    Gender
	Location  Location
	Inventory map[string]int
}

type RegistrationResult struct {
	SurvivorID string
}

type RegistrationHandler struct {
	Survivors SurvivorRepository
}

func (h *RegistrationHandler) Handle(ctx context.Context, cmd RegistrationCommand) (RegistrationResult, error) {
	s := &Survivor{
		ID:        xid.New().String(),
		Name:      cmd.Name,
		Age:       cmd.Age,
		Gender:    cmd.Gender,
		Location:  cmd.Location,
		Inventory: newResources(cmd.Inventory),
	}

	return RegistrationResult{SurvivorID: s.ID}, h.Survivors.Save(ctx, s)
}

func (h *RegistrationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var cmd RegistrationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := h.Handle(r.Context(), cmd)
	switch {
	case err == nil:
		if err := json.NewEncoder(w).Encode(res); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
