package main

import "fmt"

// Интерфейс продукта
type Shape interface {
	Draw()
}

// Конкретные продукты
type Circle struct{}

func (c *Circle) Draw() {
	fmt.Println("Рисуем круг.")
}

type Square struct{}

func (s *Square) Draw() {
	fmt.Println("Рисуем квадрат.")
}

// Фабрика
func ShapeFactory(shapeType string) Shape {
	if shapeType == "circle" {
		return &Circle{}
	} else if shapeType == "square" {
		return &Square{}
	}
	return nil
}

func main() {
	shape1 := ShapeFactory("circle")
	shape1.Draw()

	shape2 := ShapeFactory("square")
	shape2.Draw()
}
