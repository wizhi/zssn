package zssn

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InfectedQuery struct{}

type InfectedResult struct {
	PercentageInfected    float64
	PercentageNotInfected float64
}

type InfectedHandler struct {
	Conn *pgxpool.Pool
}

func (h *InfectedHandler) Handle(ctx context.Context, q InfectedQuery) (InfectedResult, error) {
	var res InfectedResult
	if h.Conn == nil {
		return res, nil
	}

	row := h.Conn.QueryRow(ctx, `
		SELECT avg(CASE WHEN flags >= 3 THEN 1 ELSE 0 END)
		FROM survivor
	`)

	var infected float64
	err := row.Scan(&infected)

	res.PercentageInfected = infected * 100
	res.PercentageNotInfected = 100 - res.PercentageInfected

	return res, err
}

func (h *InfectedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res, err := h.Handle(r.Context(), InfectedQuery{})
	switch {
	case err == nil:
		if err := json.NewEncoder(w).Encode(res); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
