package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Lualttt/lua.lt/internal/handlers"
	"github.com/Lualttt/lua.lt/internal/visits"
)

func main() {
	slog.Info("starting http server on :8080")

	go visits.Main()

	http.HandleFunc("/", handlers.Root)
	http.HandleFunc("/process", handlers.Process)
	http.HandleFunc("/visits", handlers.Visits)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	err := http.ListenAndServe(os.Getenv("LUALT_ADDRESS"), nil)
	slog.Error(err.Error())
}
