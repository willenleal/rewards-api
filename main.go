package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/willenleal/rewards-api/api"
)

func main() {
	server := api.NewServer()

	r := http.NewServeMux()

	h := api.HandlerFromMux(server, r)

	s := &http.Server{
		Handler: h,
		Addr:    ":3000",
	}

	slog.Info("Listening on port 3000")

	log.Fatal(s.ListenAndServe())
}
