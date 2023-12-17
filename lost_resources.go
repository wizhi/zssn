package zssn

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LostResourcesQuery struct{}

type LostResourcesResult struct {
	Lost int
}

type LostResourcesHandler struct {
	Conn *pgxpool.Pool
}

func (h *LostResourcesHandler) Handle(ctx context.Context, q LostResourcesQuery) (LostResourcesResult, error) {
	var res LostResourcesResult
	if h.Conn == nil {
		return res, nil
	}

	// This reports the total points lost due to _all_ infected survivors.
	// The case does specify "[..] because of infected **survivor**",
	// but I'll make an assumption that this was a typo.
	rows, err := h.Conn.Query(ctx, `
		SELECT kind, sum(quantity)
		FROM resource r
		JOIN survivor s ON s.id = r.survivor_id
		WHERE s.flags >= 3
		GROUP BY kind
	`)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			kind string
			lost int
		)
		if err := rows.Scan(&kind, &lost); err != nil {
			return res, err
		}
		res.Lost += lost * items[kind].Worth
	}

	return res, nil
}

func (h *LostResourcesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res, err := h.Handle(r.Context(), LostResourcesQuery{})
	switch {
	case err == nil:
		if err := json.NewEncoder(w).Encode(res); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
