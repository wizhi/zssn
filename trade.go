package zssn

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var (
	// ErrTradeImpossible indicates an exchange of an invalid quantity of resourecs.
	// The upcoming Scams and Thievery update will possibly remove the need for such
	// sensibilities, assuming the survivor's skill levels are high enough.
	ErrTradeImpossible = errors.New("impossible trade")
	// ErrTradeUnfair indicates an exchange of resources of unequal worth, because that'd be unfair.
	// As we all know, value is wholly objective, no matter the context. Who needs bullets anyway?
	ErrTradeUnfair = errors.New("uneven trade")
)

// Trade will exchange the offered resources from s1's inventory with the wanted resources from s2's inventory.
// If either survivor is infected, an ErrInfected error will be returned.
// If the resources being exchanged are not available in the survivors' inventories, an ErrTradeImpossible error will be returned.
// If the resources being exchanged are not of even worth, an ErrTradeUnfair error will be returned.
func Trade(s1 *Survivor, offer Resources, s2 *Survivor, want Resources) error {
	if s1.Infected() || s2.Infected() {
		return ErrInfected
	}

	if !s1.Inventory.Has(offer) || !s2.Inventory.Has(want) {
		return ErrTradeImpossible
	}

	if offer.Worth() != want.Worth() {
		return ErrTradeUnfair
	}

	for item, q := range offer {
		s1.addItem(item, -q)
		s2.addItem(item, q)
	}

	for item, q := range want {
		s1.addItem(item, q)
		s2.addItem(item, -q)
	}

	return nil
}

type TradeCommand struct {
	From, To    string
	Offer, Want map[string]int
}

type TradeHandler struct {
	Survivors SurvivorRepository
}

func (h *TradeHandler) Handle(ctx context.Context, cmd TradeCommand) error {
	s1, err := h.Survivors.Load(ctx, cmd.From)
	if err != nil {
		return err
	}

	s2, err := h.Survivors.Load(ctx, cmd.To)
	if err != nil {
		return err
	}

	offer := newResources(cmd.Offer)
	want := newResources(cmd.Want)

	if err := Trade(s1, offer, s2, want); err != nil {
		return err
	}

	return h.Survivors.Save(ctx, s1, s2)
}

func (h *TradeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var cmd TradeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cmd.From = r.Header.Get("X-Survivor")
	cmd.To = chi.URLParam(r, "survivorID")

	err := h.Handle(r.Context(), cmd)
	switch {
	case err == nil:
		w.WriteHeader(http.StatusNoContent)
	case errors.Is(err, ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
	case errors.Is(err, ErrTradeImpossible) || errors.Is(err, ErrTradeUnfair):
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
