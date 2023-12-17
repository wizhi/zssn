package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
)

func main() {
	fs := ff.NewFlagSet("zssn")

	var (
		listen = fs.StringLong("listen", ":8080", "Address to serve requests on")
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

	r := chi.NewMux()

	r.Use(
		middleware.Logger,
		middleware.Recoverer,
		middleware.AllowContentType("application/json"),
	)

	r.Get("/", health())

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
