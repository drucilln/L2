package main

import (
	"fmt"
	"testing"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"a4bc2d5e", "aaaabccddddde"},
		{"abcd", "abcd"},
		{"", ""},
		{"45", "(некорректная строка)"},
		{"qwe\\4\\5", "qwe45"},
		{"qwe\\45", "qwe44444"},
		{"qwe\\\\5", "qwe\\\\\\\\\\"},
	}

	for _, test := range tests {
		result := Unpack(test.input)
		if result != test.expected {
			fmt.Printf("Для строки %s ожидалось %s, но получено %s \n", test.input, test.expected, result)
		}

	}
}
