package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/wizhi/zssn"
	"github.com/wizhi/zssn/inmem"
	"github.com/wizhi/zssn/postgres"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
)

func main() {
	fs := ff.NewFlagSet("zssn")

	var (
		listen     = fs.StringLong("listen", ":8080", "Address to serve requests on")
		connString = fs.StringLong("postgres.conn", "", "Postgres connection string")
	)

	err := ff.Parse(fs, os.Args[1:],
		ff.WithEnvVars(),
	)
	switch {
	case errors.Is(err, ff.ErrHelp):
		fmt.Fprintf(os.Stderr, "%s\n", ffhelp.Flags(fs))
		os.Exit(0)
	case err != nil:
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	var repo zssn.SurvivorRepository = &inmem.SurvivorRepository{}
	if *connString != "" {
		pg, err := pgxpool.New(ctx, *connString)
		if err != nil {
			log.Fatalf("invalid Postgres connection string: %v", err)
		}
		repo = &postgres.SurvivorRepository{Conn: pg}
	}

	registration := &zssn.RegistrationHandler{Survivors: repo}
	status := &zssn.StatusHandler{Survivors: repo}
	checkin := &zssn.CheckInHandler{Survivors: repo}
	flag := &zssn.FlagHandler{Survivors: repo}
	trade := &zssn.TradeHandler{Survivors: repo}

	r := chi.NewMux()

	r.Use(
		middleware.Logger,
		middleware.Recoverer,
		middleware.AllowContentType("application/json"),
	)

	r.Get("/", health())

	r.Route("/survivors", func(r chi.Router) {
		r.Post("/", registration.ServeHTTP)
		r.Route("/{survivorID}", func(r chi.Router) {
			r.Get("/", status.ServeHTTP)
			r.Post("/checkins", checkin.ServeHTTP)
			r.Post("/flags", flag.ServeHTTP)
			r.Post("/trades", trade.ServeHTTP)
		})
	})

	srv := &http.Server{
		Addr:    *listen,
		Handler: r,
	}

	// TODO Add graceful shutdown of HTTP server
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server closed abruptly: %v", err)
	}
}

func health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	}
}
