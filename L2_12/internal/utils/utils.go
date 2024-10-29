package utils

import (
	"L2_12/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func ParseEvent(r *http.Request) (models.Event, error) {
	var event models.Event
	err := r.ParseForm()
	if err != nil {
		return event, fmt.Errorf("invalid form data")
	}

	userID, err := strconv.Atoi(r.Form.Get("user_id"))
	if err != nil {
		return event, fmt.Errorf("invalid user_id")
	}
	dateStr := r.Form.Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return event, fmt.Errorf("invalid date")
	}

	event.UserID = userID
	event.Title = r.Form.Get("title")
	event.Description = r.Form.Get("description")
	event.Date = date

	if idStr := r.Form.Get("event_id"); idStr != "" {
		eventID, err := strconv.Atoi(idStr)
		if err != nil {
			return event, fmt.Errorf("invalid event_id")
		}
		event.ID = eventID
	}

	return event, nil
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, errMsg string) {
	WriteJSON(w, status, map[string]string{"error": errMsg})
}

func WriteResult(w http.ResponseWriter, result interface{}) {
	WriteJSON(w, http.StatusOK, map[string]interface{}{"result": result})
}
