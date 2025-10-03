package visits

import (
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type visitsFile struct {
	Visits int `json:"visits"`
}

func Main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	SetVisits(LoadVisits())
	go autoSave()

	<-sigChan

	slog.Info("shutting down saving visits")
	SaveVisits(GetVisits())
	os.Exit(0)
}

func autoSave() {
	for {
		time.Sleep(1 * time.Hour)

		SaveVisits(GetVisits())
	}
}

func SaveVisits(amount int) {
	json, err := json.Marshal(visitsFile{Visits: amount})
	if err != nil {
		slog.Error("error marshalling visits", "error", err)
		return
	}

	err = os.WriteFile("visits.json", json, 0644)
	if err != nil {
		slog.Error("error writing visits", "error", err)
	}
}

func LoadVisits() int {
	data, err := os.ReadFile("visits.json")
	if err != nil {
		slog.Error("error loading visits", "error", err)
		data = []byte("{\"visits\": 0}")
	}

	visits := visitsFile{}
	err = json.Unmarshal(data, &visits)
	if err != nil {
		slog.Error("error unmarshalling visits", "error", err)
	}

	return visits.Visits
}
