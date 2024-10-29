package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	column  = flag.Int("k", 0, "Указание колонки для сортировки")
	numeric = flag.Bool("n", false, "Сортировать по числовому значению")
	reverse = flag.Bool("r", false, "Сортировать в обратном порядке")
	unique  = flag.Bool("u", false, "Не выводить повторяющиеся строки")
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Укажите имя входного файла")
		os.Exit(1)
	}

	fileName := flag.Arg(0)

	lines, err := readLines(fileName)
	if err != nil {
		log.Fatalf("Ошибка при чтении файла: %v", err)
	}

	lines = sortLinesFunc(lines)

	// Выводим результат
	for _, line := range lines {
		fmt.Println(line)
	}
}

func readLines(fileName string) ([]string, error) {
	var lines []string

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

type sortLines struct {
	lines  []string
	lessFn func(i, j int) bool
}

func (s sortLines) Len() int {
	return len(s.lines)
}

func (s sortLines) Swap(i, j int) {
	s.lines[i], s.lines[j] = s.lines[j], s.lines[i]
}

func (s sortLines) Less(i, j int) bool {
	return s.lessFn(i, j)
}

func getKey(line string) string {
	if *column > 0 {
		parts := strings.Fields(line)
		if len(parts) >= *column {
			return parts[*column-1]
		}
		return ""
	}
	return line
}

func compare(i, j int, lines []string) bool {
	keyI := getKey(lines[i])
	keyJ := getKey(lines[j])

	keyIEmpty := keyI == ""
	keyJEmpty := keyJ == ""

	// Обработка пустых ключей
	if keyIEmpty && keyJEmpty {
		if *reverse {
			return lines[i] > lines[j]
		}
		return lines[i] < lines[j]
	} else if keyIEmpty {
		// Строка без ключа считается больше
		return false
	} else if keyJEmpty {
		return true
	}

	// Сравнение по числовому значению
	if *numeric {
		numI, errI := strconv.ParseFloat(keyI, 64)
		numJ, errJ := strconv.ParseFloat(keyJ, 64)

		if errI == nil && errJ == nil {
			if *reverse {
				return numI > numJ
			}
			return numI < numJ
		} else if errI == nil {
			return true
		} else if errJ == nil {
			return false
		}
	}

	// Лексикографическое сравнение ключей
	if *reverse {
		return keyI > keyJ
	}
	return keyI < keyJ
}

func sortLinesFunc(lines []string) []string {
	sorted := sortLines{
		lines: lines,
		lessFn: func(i, j int) bool {
			return compare(i, j, lines)
		},
	}

	sort.Sort(sorted)

	if *unique {
		lines = removeDuplicates(sorted.lines)
	} else {
		lines = sorted.lines
	}

	return lines
}

func removeDuplicates(lines []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, line := range lines {
		if !seen[line] {
			result = append(result, line)
			seen[line] = true
		}
	}
	return result
}
