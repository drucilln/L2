package models

import (
	"fmt"
	"sync"
	"time"
)

type Event struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

type EventStore interface {
	CreateEvent(event Event) error
	UpdateEvent(event Event) error
	DeleteEvent(eventID int) error
	EventsForDay(userID int, date time.Time) ([]Event, error)
	EventsForWeek(userID int, date time.Time) ([]Event, error)
	EventsForMonth(userID int, date time.Time) ([]Event, error)
}

type InMemoryEventStore struct {
	events map[int]Event
	mu     sync.Mutex
	nextID int
}

func NewInMemoryEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{
		events: make(map[int]Event),
		nextID: 1,
	}
}

func (store *InMemoryEventStore) CreateEvent(event Event) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	event.ID = store.nextID
	store.nextID++
	store.events[event.ID] = event
	return nil
}

func (store *InMemoryEventStore) UpdateEvent(event Event) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	if _, exists := store.events[event.ID]; !exists {
		return fmt.Errorf("event not found")
	}
	store.events[event.ID] = event
	return nil
}

func (store *InMemoryEventStore) DeleteEvent(eventID int) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	if _, exists := store.events[eventID]; !exists {
		return fmt.Errorf("event not found")
	}
	delete(store.events, eventID)
	return nil
}

func (store *InMemoryEventStore) EventsForDay(userID int, date time.Time) ([]Event, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	var events []Event
	for _, event := range store.events {
		if event.UserID == userID && sameDay(event.Date, date) {
			events = append(events, event)
		}
	}
	return events, nil
}

func (store *InMemoryEventStore) EventsForWeek(userID int, date time.Time) ([]Event, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	var events []Event
	weekStart := date.Truncate(24 * time.Hour)
	weekEnd := weekStart.AddDate(0, 0, 7)
	for _, event := range store.events {
		if event.UserID == userID && event.Date.After(weekStart) && event.Date.Before(weekEnd) {
			events = append(events, event)
		}
	}
	return events, nil
}

func (store *InMemoryEventStore) EventsForMonth(userID int, date time.Time) ([]Event, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	var events []Event
	monthStart := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	monthEnd := monthStart.AddDate(0, 1, 0)
	for _, event := range store.events {
		if event.UserID == userID && event.Date.After(monthStart) && event.Date.Before(monthEnd) {
			events = append(events, event)
		}
	}
	return events, nil
}

func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}
