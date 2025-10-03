package main

import (
	"log/slog"
	"net/http"

	"github.com/Lualttt/lua.lt/internal/handlers"
	"github.com/Lualttt/lua.lt/internal/visits"
)

func main() {
	slog.Info("starting http server on :8080")

	go visits.Main()

	http.HandleFunc("/", handlers.Root)
	http.HandleFunc("/visits", handlers.Visits)

	err := http.ListenAndServe(":8080", nil)
	slog.Error(err.Error())
}
