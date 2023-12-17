package zssn_test

import (
	"context"
	"testing"

	"github.com/wizhi/zssn"
	"github.com/wizhi/zssn/inmem"

	"github.com/google/go-cmp/cmp"
)

func TestRegistrationHandler_Handle(t *testing.T) {
	repo := &inmem.SurvivorRepository{}
	registration := &zssn.RegistrationHandler{Survivors: repo}
	status := &zssn.StatusHandler{Survivors: repo}

	ctx := context.Background()
	cmd := zssn.RegistrationCommand{
		Name:     "Simon",
		Age:      29,
		Gender:   zssn.Male,
		Location: zssn.LatLong(55.693573397195664, 12.569842474809072),
		Inventory: map[string]int{
			zssn.Water.Kind:      3, // Soda counts, right?
			zssn.Food.Kind:       2, // Chips too, I hope.
			zssn.Medication.Kind: 1, // Need to restock on painkillers sometime soon.
		},
	}

	res, err := registration.Handle(ctx, cmd)
	if err != nil {
		t.Fatalf("registration failed: %v", err)
	}

	got, err := status.Handle(ctx, zssn.StatusQuery{SurvivorID: res.SurvivorID})
	if err != nil {
		t.Fatalf("status failed: %v", err)
	}

	want := zssn.StatusResult{
		Name:      cmd.Name,
		Age:       cmd.Age,
		Gender:    cmd.Gender,
		Location:  cmd.Location,
		Inventory: cmd.Inventory,
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("status does not match registration: (-want +got)\n%s", diff)
	}
}
