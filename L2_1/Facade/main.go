package main

import "fmt"

type Water struct{}

func (w *Water) Boil() {
	fmt.Println("Кипячение воды...")
}

type TeaLeaf struct{}

func (t *TeaLeaf) Add() {
	fmt.Println("Добавление чайных листьев...")
}

type Sugar struct{}

func (s *Sugar) Add() {
	fmt.Println("Добавление сахара...")
}

// Фасад
type TeaFacade struct {
	water   *Water
	teaLeaf *TeaLeaf
	sugar   *Sugar
}

func NewTeaFacade() *TeaFacade {
	return &TeaFacade{
		water:   &Water{},
		teaLeaf: &TeaLeaf{},
		sugar:   &Sugar{},
	}
}

func (t *TeaFacade) MakeTea() {
	t.water.Boil()
	t.teaLeaf.Add()
	t.sugar.Add()
	fmt.Println("Чай готов!")
}

func main() {
	teaMaker := NewTeaFacade()
	teaMaker.MakeTea()
}
