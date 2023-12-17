package zssn

import (
	"context"
	"errors"
)

var (
	// ErrInfected is returned when survivors attempt to perform actions infected aren't capable of.
	ErrInfected = errors.New("survivor is infected")

	// ErrNotFound indicates than an entity does not exist in the system.
	ErrNotFound = errors.New("not found")
)

// A Survivor is your everyday human being, just trying to get by.
// These may or may not be capable of thought. Oh, and some are zombies too, I guess.
type Survivor struct {
	ID string

	Name      string
	Age       int
	Gender    Gender
	Location  Location
	Flags     int
	Inventory Resources
}

type SurvivorRepository interface {
	// Load will find any persisted survivor with the given identifier.
	// If no survivor is found, ErrNotFound will be returned.
	Load(ctx context.Context, id string) (*Survivor, error)

	// Save will persist the given survivor(s) atomically.
	Save(ctx context.Context, ss ...*Survivor) error
}

// Gender is complicated.
type Gender string

const (
	Male   Gender = "M"
	Female Gender = "F"
)

const infectedFlagThreshold = 3

// Flag will report the survivor as being infected.
// Once a survivor has been flagged enough times, they will become infected.
func (s *Survivor) Flag() {
	s.Flags++
}

// Infected indicates whether the survivor is considered to be infected.
func (s *Survivor) Infected() bool {
	return s.Flags >= infectedFlagThreshold
}

// A Location is a set of coordinates.
type Location struct {
	Latitude, Longitude float64
}

// LatLong returns the location of the given coordinates.
// This is just intended as a convenience constructor for laziness.
func LatLong(lat, long float64) Location {
	return Location{
		Latitude:  lat,
		Longitude: long,
	}
}

// CheckIn will register the given coordinates as the survivor's last location.
func (s *Survivor) CheckIn(lat, long float64) {
	s.Location = LatLong(lat, long)
}

// Resources is a set of items in various quantities.
// The default value is considered an empty set of items.
type Resources map[Item]int

// An Item describes a piece of property of some worth.
type Item struct {
	Kind  string
	Worth int
}

// This is populated at program startup, and should not be modified.
// Realistically, we'd keep a persistent register somewhere else,
// or at least do some automated generation, but you know, KISS.
var items map[string]Item

var (
	Ammunition = Item{"ammunition", 1}
	Food       = Item{"food", 3}
	Medication = Item{"medication", 2}
	Water      = Item{"water", 4}
)

func newResources(m map[string]int) Resources {
	r := make(Resources, len(m))
	for kind, q := range m {
		if item, ok := items[kind]; ok {
			r[item] = q
		}
	}
	return r
}

// Has checks if given resources are a subset of the current resources.
// This will always be false, should any of the sub resources be of a negative quantity.
func (m Resources) Has(sub Resources) bool {
	if len(m) == 0 {
		return len(sub) == 0
	}
	for item, q := range sub {
		if q < 0 || m[item] < q {
			return false
		}
	}
	return true
}

// Worth determines the total worth of the current resources.
func (m Resources) Worth() int {
	var n int
	for item, q := range m {
		n += item.Worth * q
	}
	return n
}

// addItem adds a quantity of the given item to the survivor's inventory.
func (s *Survivor) addItem(item Item, quantity int) error {
	if s.Inventory == nil {
		s.Inventory = make(Resources)
	}

	if s.Infected() {
		return ErrInfected
	}

	n := s.Inventory[item] + quantity
	switch {
	case n > 0:
		s.Inventory[item] = n
	case n == 0:
		delete(s.Inventory, item)
	default:
		return errors.New("not enough resources")
	}

	return nil
}

func init() {
	items = map[string]Item{
		Ammunition.Kind: Ammunition,
		Food.Kind:       Food,
		Medication.Kind: Medication,
		Water.Kind:      Water,
	}
}
