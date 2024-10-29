package main

import "fmt"

// Интерфейс стратегии
type PaymentStrategy interface {
	Pay(amount float64)
}

// Конкретные стратегии
type CreditCard struct {
	Name   string
	CardNo string
}

func (c *CreditCard) Pay(amount float64) {
	fmt.Printf("Оплачено %.2f с помощью кредитной карты %s\n", amount, c.CardNo)
}

type PayPal struct {
	Email string
}

func (p *PayPal) Pay(amount float64) {
	fmt.Printf("Оплачено %.2f с помощью PayPal аккаунта %s\n", amount, p.Email)
}

// Контекст
type Order struct {
	amount   float64
	strategy PaymentStrategy
}

func (o *Order) SetPaymentStrategy(s PaymentStrategy) {
	o.strategy = s
}

func (o *Order) Pay() {
	o.strategy.Pay(o.amount)
}

func main() {
	order := &Order{amount: 100.0}

	// Клиент выбирает оплату кредитной картой
	order.SetPaymentStrategy(&CreditCard{Name: "Иван Иванов", CardNo: "1234-5678-9012-3456"})
	order.Pay()

	// Клиент меняет решение и выбирает PayPal
	order.SetPaymentStrategy(&PayPal{Email: "ivanov@example.com"})
	order.Pay()
}
