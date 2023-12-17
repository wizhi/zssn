package postgres

import (
	"context"
	"fmt"

	"github.com/wizhi/zssn"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SurvivorRepository struct {
	Conn *pgxpool.Pool
}

func (r *SurvivorRepository) Load(ctx context.Context, id string) (*zssn.Survivor, error) {
	return nil, nil
}

func (r *SurvivorRepository) Save(ctx context.Context, ss ...*zssn.Survivor) error {
	tx, err := r.Conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, s := range ss {
		loc := pgtype.Point{
			P:     pgtype.Vec2{X: s.Location.Latitude, Y: s.Location.Longitude},
			Valid: true,
		}

		if _, err := tx.Exec(ctx, `
			INSERT INTO survivor (
				id,
				name,
				gender,
				location,
				flags
			)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (id) DO UPDATE SET
				name = $2,
				gender = $3,
				location = $4,
				flags = $5
		`, s.ID, s.Name, s.Gender, loc, s.Flags); err != nil {
			return fmt.Errorf("failed to upsert survivor: %v", err)
		}

		if _, err := tx.Exec(ctx, "DELETE FROM resource WHERE survivor_id = $1", s.ID); err != nil {
			return fmt.Errorf("failed to clear resources: %v", err)
		}

		for item, q := range s.Inventory {
			// TODO Roll into singular insertion for all resources
			if _, err := tx.Exec(ctx, `
				INSERT INTO resource (survivor_id, kind, quantity)
				VALUES ($1, $2, $3)
			`, s.ID, item.Kind, q); err != nil {
				return fmt.Errorf("failed to insert resource: %v", err)
			}
		}
	}

	return tx.Commit(ctx)
}
