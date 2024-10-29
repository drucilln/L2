package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	tests := []string{
		"a4bc2d5e",
		"abcd",
		"45",
		"",
		"qwe\\4\\5",
		"qwe\\45",
		"qwe\\\\5",
	}

	for _, testStr := range tests {
		res := Unpack(testStr)
		fmt.Println(res)
	}

}

func Unpack(str string) string {
	buider := strings.Builder{}
	var prevRune rune
	var isEscaped bool
	for _, v := range str {

		if isEscaped {
			// Добавляем символ как есть
			buider.WriteRune(v)
			prevRune = v
			isEscaped = false
			continue
		}

		if v == '\\' {
			// Включаем режим экранирования
			isEscaped = true
			continue
		}
		if unicode.IsDigit(v) && prevRune != 0 {
			r, _ := strconv.Atoi(string(v))
			for ; r > 1; r-- {
				buider.WriteRune(prevRune)
			}

		} else if unicode.IsDigit(v) && prevRune == 0 {
			return "(некорректная строка)"
		} else if unicode.IsLetter(v) {
			buider.WriteRune(v)
			prevRune = v
		}
	}
	return buider.String()
}
