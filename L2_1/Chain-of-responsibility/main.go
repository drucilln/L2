package main

import "fmt"

// Интерфейс обработчика
type Handler interface {
	SetNext(Handler)
	Handle(email *Email)
}

// Структура Email
type Email struct {
	Content  string
	Spam     bool
	Phishing bool
	Virus    bool
}

// Базовый обработчик
type BaseHandler struct {
	next Handler
}

func (h *BaseHandler) SetNext(next Handler) {
	h.next = next
}

func (h *BaseHandler) Handle(email *Email) {
	if h.next != nil {
		h.next.Handle(email)
	}
}

// Конкретные обработчики
type SpamHandler struct {
	BaseHandler
}

func (h *SpamHandler) Handle(email *Email) {
	if email.Spam {
		fmt.Println("Email помечен как спам.")
		return
	}
	h.BaseHandler.Handle(email)
}

type PhishingHandler struct {
	BaseHandler
}

func (h *PhishingHandler) Handle(email *Email) {
	if email.Phishing {
		fmt.Println("Email содержит фишинг.")
		return
	}
	h.BaseHandler.Handle(email)
}

type VirusHandler struct {
	BaseHandler
}

func (h *VirusHandler) Handle(email *Email) {
	if email.Virus {
		fmt.Println("Email содержит вирус.")
		return
	}
	h.BaseHandler.Handle(email)
}

type InboxHandler struct {
	BaseHandler
}

func (h *InboxHandler) Handle(email *Email) {
	fmt.Println("Email доставлен в папку 'Входящие'.")
}

func main() {
	email := &Email{Content: "Привет!", Spam: false, Phishing: false, Virus: false}

	spamHandler := &SpamHandler{}
	phishingHandler := &PhishingHandler{}
	virusHandler := &VirusHandler{}
	inboxHandler := &InboxHandler{}

	spamHandler.SetNext(phishingHandler)
	phishingHandler.SetNext(virusHandler)
	virusHandler.SetNext(inboxHandler)

	spamHandler.Handle(email)
}
