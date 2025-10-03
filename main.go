package main

import (
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"unicode/utf8"
)

type PageVariables struct {
	Views int
}

func (pageVariables PageVariables) ViewsAsArray() []string {
	var text string = strconv.Itoa(pageVariables.Views)
	characters := make([]string, utf8.RuneCountInString(text))

	for i, char := range text {
		characters[i] = string(char)
	}

	return characters
}

func main() {
	slog.Info("starting http server on :8080")

	http.HandleFunc("/", index)

	err := http.ListenAndServe(":8080", nil)
	slog.Error(err.Error())
}

func index(w http.ResponseWriter, r *http.Request) {
	pageVariables := PageVariables{
		Views: 3444,
	}

	t, err := template.ParseFiles("./web/templates/index.html")
	if err != nil {
		http.Error(w, "Internal server error!", http.StatusInternalServerError)
		slog.Error("index page error", "error", err.Error())
		return
	}

	err = t.Execute(w, pageVariables)
	if err != nil {
		http.Error(w, "Internal server error!", http.StatusInternalServerError)
		slog.Error("index page error", "error", err.Error())
		return
	}

	slog.Info("index page loaded")
}
