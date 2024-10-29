package main

import "fmt"

// Интерфейс состояния
type State interface {
	Next(*TrafficLight)
	Show()
}

// Конкретные состояния
type RedState struct{}

func (s *RedState) Next(t *TrafficLight) {
	t.SetState(&GreenState{})
}

func (s *RedState) Show() {
	fmt.Println("Светофор красный. Стойте.")
}

type YellowState struct{}

func (s *YellowState) Next(t *TrafficLight) {
	t.SetState(&RedState{})
}

func (s *YellowState) Show() {
	fmt.Println("Светофор жёлтый. Приготовьтесь.")
}

type GreenState struct{}

func (s *GreenState) Next(t *TrafficLight) {
	t.SetState(&YellowState{})
}

func (s *GreenState) Show() {
	fmt.Println("Светофор зелёный. Можно идти.")
}

// Контекст
type TrafficLight struct {
	state State
}

func (t *TrafficLight) SetState(s State) {
	t.state = s
}

func (t *TrafficLight) Change() {
	t.state.Next(t)
}

func (t *TrafficLight) Show() {
	t.state.Show()
}

func main() {
	light := &TrafficLight{state: &RedState{}}

	for i := 0; i < 5; i++ {
		light.Show()
		light.Change()
	}
}
