package main

import "fmt"

// Продукт
type Message struct {
	greeting   string
	body       string
	conclusion string
}

func (m *Message) Display() {
	fmt.Println(m.greeting)
	fmt.Println(m.body)
	fmt.Println(m.conclusion)
}

// Строитель
type MessageBuilder struct {
	message *Message
}

func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{&Message{}}
}

func (b *MessageBuilder) SetGreeting(greeting string) *MessageBuilder {
	b.message.greeting = greeting
	return b
}

func (b *MessageBuilder) SetBody(body string) *MessageBuilder {
	b.message.body = body
	return b
}

func (b *MessageBuilder) SetConclusion(conclusion string) *MessageBuilder {
	b.message.conclusion = conclusion
	return b
}

func (b *MessageBuilder) Build() *Message {
	return b.message
}

func main() {
	builder := NewMessageBuilder()
	message := builder.
		SetGreeting("Здравствуйте!").
		SetBody("Поздравляем вас с днём рождения!").
		SetConclusion("С наилучшими пожеланиями, команда XYZ.").
		Build()
	message.Display()
}
