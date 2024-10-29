package handlers

import (
	"L2_12/internal/models"
	"L2_12/internal/utils"
	"net/http"
	"strconv"
	"time"
)

func CreateEventHandler(store models.EventStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		event, err := utils.ParseEvent(r)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		err = store.CreateEvent(event)
		if err != nil {
			utils.WriteError(w, http.StatusServiceUnavailable, err.Error())
			return
		}

		utils.WriteResult(w, "Event created")
	}
}

func UpdateEventHandler(store models.EventStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		event, err := utils.ParseEvent(r)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		err = store.UpdateEvent(event)
		if err != nil {
			utils.WriteError(w, http.StatusServiceUnavailable, err.Error())
			return
		}

		utils.WriteResult(w, "Event updated")
	}
}

func DeleteEventHandler(store models.EventStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		err := r.ParseForm()
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Invalid form data")
			return
		}

		eventIDStr := r.Form.Get("event_id")
		eventID, err := strconv.Atoi(eventIDStr)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Invalid event_id")
			return
		}

		err = store.DeleteEvent(eventID)
		if err != nil {
			utils.WriteError(w, http.StatusServiceUnavailable, err.Error())
			return
		}

		utils.WriteResult(w, "Event deleted")
	}
}

func EventsHandler(store models.EventStore, period string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		userIDStr := r.URL.Query().Get("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Invalid user_id")
			return
		}

		dateStr := r.URL.Query().Get("date")
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Invalid date")
			return
		}

		var events []models.Event
		switch period {
		case "day":
			events, err = store.EventsForDay(userID, date)
		case "week":
			events, err = store.EventsForWeek(userID, date)
		case "month":
			events, err = store.EventsForMonth(userID, date)
		default:
			utils.WriteError(w, http.StatusInternalServerError, "Invalid period")
			return
		}

		if err != nil {
			utils.WriteError(w, http.StatusServiceUnavailable, err.Error())
			return
		}

		utils.WriteResult(w, events)
	}
}
