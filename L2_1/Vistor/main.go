package main

import "fmt"

// Интерфейс элемента
type Document interface {
	Accept(Visitor)
}

// Конкретные элементы
type TextDocument struct {
	Text string
}

func (d *TextDocument) Accept(v Visitor) {
	v.VisitTextDocument(d)
}

type ImageDocument struct {
	ImageData []byte
}

func (d *ImageDocument) Accept(v Visitor) {
	v.VisitImageDocument(d)
}

// Интерфейс посетителя
type Visitor interface {
	VisitTextDocument(*TextDocument)
	VisitImageDocument(*ImageDocument)
}

// Конкретный посетитель
type PrintVisitor struct{}

func (p *PrintVisitor) VisitTextDocument(d *TextDocument) {
	fmt.Println("Печать текстового документа:")
	fmt.Println(d.Text)
}

func (p *PrintVisitor) VisitImageDocument(d *ImageDocument) {
	fmt.Println("Печать изображения:")
	fmt.Println("[Изображение]")
}

func main() {
	documents := []Document{
		&TextDocument{Text: "Привет, мир!"},
		&ImageDocument{ImageData: []byte{0xFF, 0xD8, 0xFF}},
	}

	printer := &PrintVisitor{}

	for _, doc := range documents {
		doc.Accept(printer)
	}
}
