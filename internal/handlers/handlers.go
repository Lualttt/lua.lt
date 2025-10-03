package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"text/template"
	"unicode/utf8"

	"github.com/Lualttt/lua.lt/internal/visits"
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

func Root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		Index(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	visits.SetVisits(visits.GetVisits() + 1)

	pageVariables := PageVariables{
		Views: visits.GetVisits(),
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

func Visits(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, visits.GetVisits())
}
