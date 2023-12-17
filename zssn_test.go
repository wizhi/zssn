package zssn

import (
	"fmt"
	"strings"
	"testing"
)

func TestSurvivor_Infected(t *testing.T) {
	tests := map[string]struct {
		start, end int
		want       bool
	}{
		"under threshold": {0, infectedFlagThreshold - 1, false},
		// We assume more than 10k flags also works, for the sake of keeping test times reasonable.
		"over threshold": {infectedFlagThreshold, 10000, true},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			s := new(Survivor)
			for i := 0; i < tt.start; i++ {
				s.Flag()
			}

			for i := tt.start; i < tt.end; i++ {
				s.Flag()

				if got := s.Infected(); got != tt.want {
					t.Fatalf("(Flags = %d).Infected() = %t, want %t", s.Flags, got, tt.want)
				}
			}
		})
	}
}

func TestResources_Has(t *testing.T) {
	tests := map[string]struct {
		resources, subset Resources
		want              bool
	}{
		"zero":     {nil, nil, true},
		"equal":    {Resources{Water: 1}, Resources{Water: 1}, true},
		"more":     {Resources{Water: 2}, Resources{Water: 1}, true},
		"less":     {Resources{Water: 1}, Resources{Water: 2}, false},
		"negative": {Resources{Water: 1}, Resources{Water: -1}, false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.resources.Has(tt.subset)
			if got != tt.want {
				t.Errorf("(%s).Has(%s) = %t, want %t",
					resources(tt.resources), resources(tt.subset),
					got, tt.want,
				)
			}
		})
	}
}

func TestResources_Worth(t *testing.T) {
	tests := map[string]struct {
		resources Resources
		want      int
	}{
		"zero":    {nil, 0},
		"empty":   {Resources{}, 0},
		"single":  {Resources{Water: 1}, Water.Worth},
		"product": {Resources{Water: 2}, Water.Worth * 2},
		"sum":     {Resources{Water: 1, Food: 1}, Water.Worth + Food.Worth},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.resources.Worth()
			if got != tt.want {
				t.Errorf("(%s).Worth() = %d, want %d",
					resources(tt.resources),
					got, tt.want,
				)
			}
		})
	}
}

func resources(r Resources) string {
	s := make([]string, 0, len(r))
	for item, q := range r {
		s = append(s, fmt.Sprintf("%d %s", q, item.Kind))
	}
	return fmt.Sprintf("[%s]", strings.Join(s, "; "))
}
