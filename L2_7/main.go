package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	fields    = flag.String("f", "", "выбрать поля (колонки)")
	delimiter = flag.String("d", "\t", "использовать другой разделитель")
	separated = flag.Bool("s", false, "только строки с разделителем")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if *fields == "" {
		fmt.Println("Необходимо указать поля с помощью опции -f")
		flag.Usage()
		os.Exit(1)
	}
}

func parseFields(fieldsStr string) ([]int, error) {
	var fields []int
	ranges := strings.Split(fieldsStr, ",")
	for _, r := range ranges {
		if strings.Contains(r, "-") {
			bounds := strings.Split(r, "-")
			if len(bounds) != 2 {
				return nil, fmt.Errorf("неверный формат диапазона: %s", r)
			}
			start, err1 := strconv.Atoi(bounds[0])
			end, err2 := strconv.Atoi(bounds[1])
			if err1 != nil || err2 != nil || start > end {
				return nil, fmt.Errorf("неверный диапазон: %s", r)
			}
			for i := start; i <= end; i++ {
				fields = append(fields, i)
			}
		} else {
			fieldNum, err := strconv.Atoi(r)
			if err != nil {
				return nil, fmt.Errorf("неверный номер поля: %s", r)
			}
			fields = append(fields, fieldNum)
		}
	}
	return fields, nil
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

func processLines(lines []string, fields []int, delimiter string, separated bool) []string {
	var result []string
	for _, line := range lines {
		if separated && !strings.Contains(line, delimiter) {
			continue
		}
		tokens := strings.Split(line, delimiter)
		var selected []string
		for _, fieldNum := range fields {
			if fieldNum-1 < len(tokens) && fieldNum > 0 {
				selected = append(selected, tokens[fieldNum-1])
			}
		}
		result = append(result, strings.Join(selected, delimiter))
	}
	return result
}

func printResult(lines []string) {
	for _, line := range lines {
		fmt.Println(line)
	}
}

func main() {
	// Парсим флаги (флаги уже обработаны в init)

	// Парсим поля
	fieldsNums, err := parseFields(*fields)
	if err != nil {
		fmt.Printf("Ошибка при разборе полей: %v\n", err)
		os.Exit(1)
	}

	// Читаем входные данные
	lines, err := readInput()
	if err != nil {
		fmt.Printf("Ошибка при чтении входных данных: %v\n", err)
		os.Exit(1)
	}

	// Обрабатываем строки
	result := processLines(lines, fieldsNums, *delimiter, *separated)

	// Выводим результат
	printResult(result)
}
