package main

import (
	"L2_12/internal/handlers"
	"L2_12/internal/middleware"
	"L2_12/internal/models"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := ":8080"

	// Инициализация хранилища событий
	store := models.NewInMemoryEventStore()

	// Создание HTTP маршрутов и обработчиков
	mux := http.NewServeMux()
	mux.Handle("/create_event", handlers.CreateEventHandler(store))
	mux.Handle("/update_event", handlers.UpdateEventHandler(store))
	mux.Handle("/delete_event", handlers.DeleteEventHandler(store))
	mux.Handle("/events_for_day", handlers.EventsHandler(store, "day"))
	mux.Handle("/events_for_week", handlers.EventsHandler(store, "week"))
	mux.Handle("/events_for_month", handlers.EventsHandler(store, "month"))

	// Применение middleware для логирования запросов
	loggedMux := middleware.LoggingMiddleware(mux)

	fmt.Printf("Сервер запущен на порту %s\n", port)
	err := http.ListenAndServe(port, loggedMux)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
