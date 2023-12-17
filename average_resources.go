package zssn

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AverageResourcesQuery struct{}

type AverageResourcesResult map[string]int

type AverageResourcesHandler struct {
	Conn *pgxpool.Pool
}

func (h *AverageResourcesHandler) Handle(ctx context.Context, q AverageResourcesQuery) (AverageResourcesResult, error) {
	if h.Conn == nil {
		return nil, nil
	}

	rows, err := h.Conn.Query(ctx, `
		WITH
			uninfected AS (
				SELECT *
				FROM survivor
				WHERE flags < 3
			),
			resources AS (
			    SELECT kind, quantity
			    FROM resource r
			    JOIN uninfected ON uninfected.id = r.survivor_id
			)
		SELECT kind, sum(quantity) / (SELECT count(*) FROM uninfected)
		FROM resources
		GROUP BY kind
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make(AverageResourcesResult)
	for rows.Next() {
		var (
			kind string
			avg  int
		)
		if err := rows.Scan(&kind, &avg); err != nil {
			return nil, err
		}
		res[kind] = avg
	}

	return res, nil
}

func (h *AverageResourcesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res, err := h.Handle(r.Context(), AverageResourcesQuery{})
	switch {
	case err == nil:
		if err := json.NewEncoder(w).Encode(res); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
