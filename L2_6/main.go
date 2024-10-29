package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

var (
	after      = flag.Int("A", 0, "печать +N строк после совпадения")
	before     = flag.Int("B", 0, "печать +N строк до совпадения")
	context    = flag.Int("C", 0, "печать ±N строк вокруг совпадения")
	count      = flag.Bool("c", false, "количество строк")
	ignoreCase = flag.Bool("i", false, "игнорировать регистр")
	invert     = flag.Bool("v", false, "вместо совпадения, исключать")
	fixed      = flag.Bool("F", false, "точное совпадение со строкой, не паттерн")
	lineNum    = flag.Bool("n", false, "напечатать номер строки")
)

type line struct {
	number  int    // Номер строки
	text    string // Текст строки
	matched bool   // Флаг совпадения
}

func main() {
	flag.Parse()
	adjustContextFlags()

	if flag.NArg() < 1 {
		fmt.Println("Usage: grep [options] pattern")
		os.Exit(1)
	}

	pattern := flag.Arg(0)
	
	lines, err := readInput()
	if err != nil {
		log.Fatalf("Ошибка при чтении входных данных: %v", err)
	}

	result := grep(pattern, lines)

	printResult(result)
}

func adjustContextFlags() {
	if *context > 0 {
		*after = *context
		*before = *context
	}
}

func readInput() ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

// Функция для проверки совпадения строки с паттерном
func isMatch(pattern, text string) bool {
	if *ignoreCase {
		pattern = strings.ToLower(pattern)
		text = strings.ToLower(text)
	}

	if *fixed {
		return text == pattern
	} else {
		matched, err := regexp.MatchString(pattern, text)
		if err != nil {
			log.Fatalf("Ошибка в регулярном выражении: %v", err)
		}
		return matched
	}
}

// Функция для удаления дубликатов строк (по номеру строки)
func removeDuplicateLines(lines []line) []line {
	seen := make(map[int]struct{})
	var result []line
	for _, l := range lines {
		if _, ok := seen[l.number]; !ok {
			seen[l.number] = struct{}{}
			result = append(result, l)
		}
	}
	return result
}

// Основная функция grep
func grep(pattern string, lines []string) []line {
	var result []line
	lineMatches := make([]bool, len(lines))

	for i, text := range lines {
		matched := isMatch(pattern, text)
		if *invert {
			matched = !matched
		}
		lineMatches[i] = matched
	}

	// Второй проход: собираем строки с учетом контекста
	for i := range lines {
		matched := lineMatches[i]
		if matched {
			start := i - *before
			if start < 0 {
				start = 0
			}
			end := i + *after
			if end >= len(lines) {
				end = len(lines) - 1
			}
			for j := start; j <= end; j++ {
				l := line{
					number:  j + 1,
					text:    lines[j],
					matched: lineMatches[j],
				}
				result = append(result, l)
			}
		}
	}

	result = removeDuplicateLines(result)

	sort.Slice(result, func(i, j int) bool {
		return result[i].number < result[j].number
	})

	return result
}

func printResult(lines []line) {
	if *count {
		// Если задан флаг -c, выводим количество совпадений
		count := 0
		for _, l := range lines {
			if l.matched {
				count++
			}
		}
		fmt.Println(count)
		return
	}

	for _, l := range lines {
		// Если задан флаг -n, выводим номер строки
		if *lineNum {
			fmt.Printf("%d:", l.number)
		}
		fmt.Println(l.text)
	}
}
