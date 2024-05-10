package main

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"io"
	"log"
	"net/http"
)

type Schedule struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Class string `json:"class"`
}

var schedule []Schedule

func main() {
	http.HandleFunc("/createEvents", handleJson)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleJson(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
	err = json.Unmarshal(body, &schedule)
	if err != nil {
		log.Fatalf("Error parsing JSON: %s", err)
	}
	createEvent(schedule)
}

func createEvent(schedule []Schedule) {
	calendarID := ""
	keyFile := "credentials.json"

	ctx := context.Background()
	srv, err := calendar.NewService(ctx, option.WithCredentialsFile(keyFile))
	if err != nil {
		log.Fatalf("Failed to create Calendar service: %v", err)
	}
	for _, item := range schedule {
		if item.Class != "" {
			event := &calendar.Event{
				Summary: item.Class,
				Start: &calendar.EventDateTime{
					DateTime: item.Start,
					TimeZone: "Europe/Samara",
				},
				End: &calendar.EventDateTime{
					DateTime: item.End,
					TimeZone: "Europe/Samara",
				},
			}
			event, err = srv.Events.Insert(calendarID, event).Do()
			if err != nil {
				log.Fatalf("Error creating event: %v", err)
			}
			fmt.Printf("Event created: %s\n", event.HtmlLink)
		}
	}
}
