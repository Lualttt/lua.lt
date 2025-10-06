package main

import (
	"log/slog"
	"net/http"

	"github.com/Lualttt/lua.lt/internal/handlers"
	"github.com/Lualttt/lua.lt/internal/visits"
	"github.com/Lualttt/lua.lt/web"
)

func main() {
	slog.Info("starting http server on 0.0.0.0:8080")

	go visits.Main()

	http.HandleFunc("GET /{$}", handlers.Index)
	http.HandleFunc("POST /process", handlers.Process)
	http.HandleFunc("GET /visits", handlers.Visits)
	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(web.StaticContent))))

	err := http.ListenAndServe("0.0.0.0:8080", nil)
	slog.Error(err.Error())
}
