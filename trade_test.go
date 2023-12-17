package zssn

import (
	"errors"
	"maps"
	"testing"
)

type trade struct {
	survivor *Survivor
	exchange Resources
	after    Resources
}

func (t trade) clone() trade {
	s := *t.survivor
	s.Inventory = maps.Clone(s.Inventory)

	return trade{
		survivor: &s,
		exchange: maps.Clone(t.exchange),
		after:    maps.Clone(t.after),
	}
}

func TestTrade(t *testing.T) {
	type test struct {
		s1  trade
		s2  trade
		err error
	}

	tests := map[string]test{
		"zero": {
			s1: trade{survivor: &Survivor{}},
			s2: trade{survivor: &Survivor{}},
		},
		"swap": {
			s1: trade{
				survivor: &Survivor{Inventory: Resources{Water: 1}},
				exchange: Resources{Water: 1},
				after:    Resources{Medication: 2},
			},
			s2: trade{
				survivor: &Survivor{Inventory: Resources{Medication: 2}},
				exchange: Resources{Medication: 2},
				after:    Resources{Water: 1},
			},
		},
		"add": {
			s1: trade{
				survivor: &Survivor{Inventory: Resources{Water: 1, Ammunition: 4}},
				exchange: Resources{Ammunition: 4},
				after:    Resources{Water: 2},
			},
			s2: trade{
				survivor: &Survivor{Inventory: Resources{Water: 2}},
				exchange: Resources{Water: 1},
				after:    Resources{Water: 1, Ammunition: 4},
			},
		},
		"infected": {
			s1: trade{
				survivor: &Survivor{Flags: infectedFlagThreshold},
			},
			s2: trade{
				survivor: &Survivor{},
			},
			err: ErrInfected,
		},
		"unavailable": {
			s1: trade{
				survivor: &Survivor{Inventory: Resources{Water: 1}},
				exchange: Resources{Ammunition: 1},
			},
			s2: trade{
				survivor: &Survivor{Inventory: Resources{Medication: 2}},
				exchange: Resources{Medication: 2},
			},
			err: ErrTradeImpossible,
		},
		"negative": {
			s1: trade{
				survivor: &Survivor{Inventory: Resources{Water: -1}},
				exchange: Resources{Ammunition: 1},
			},
			s2: trade{
				survivor: &Survivor{Inventory: Resources{Medication: 2}},
				exchange: Resources{Medication: 2},
			},
			err: ErrTradeImpossible,
		},
		"unfair": {
			s1: trade{
				survivor: &Survivor{Inventory: Resources{Ammunition: 1}},
				exchange: Resources{Ammunition: 1},
			},
			s2: trade{
				survivor: &Survivor{Inventory: Resources{Water: 1}},
				exchange: Resources{Water: 1},
			},
			err: ErrTradeUnfair,
		},
	}

	for name, tt := range tests {
		tests[name+"_rev"] = test{
			s1:  tt.s2.clone(),
			s2:  tt.s1.clone(),
			err: tt.err,
		}
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := Trade(tt.s1.survivor, tt.s1.exchange, tt.s2.survivor, tt.s2.exchange)
			if tt.err != nil {
				if !errors.Is(err, tt.err) {
					t.Errorf("Trade(%v, %v, %v, %v) = %v, want %v",
						tt.s1.survivor, tt.s1.exchange,
						tt.s2.survivor, tt.s2.exchange,
						err, tt.err,
					)
				}
			} else {
				if !maps.Equal(tt.s1.survivor.Inventory, tt.s1.after) {
					t.Errorf("s1's inventory = %+v, want %+v",
						resources(tt.s1.survivor.Inventory), resources(tt.s1.after),
					)
				}
				if !maps.Equal(tt.s2.survivor.Inventory, tt.s2.after) {
					t.Errorf("s2's inventory = %+v, want %+v",
						resources(tt.s2.survivor.Inventory), resources(tt.s2.after),
					)
				}
			}
		})
	}
}
