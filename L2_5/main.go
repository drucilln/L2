package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	anagrm := []string{
		"пятка", "пятак", "тяпка", "листок", "слиток", "столик",
		"кот", "ток", "кто", "окт", "крот", "торк",
		"горка", "грока", "рогак", "агрок",
		"Кот", "кТо", "Ток",
		"дом", "мод", "дым",
		"кошка", "собака",
	}
	x := findAnagrams(anagrm)
	//for i, v := range x {
	//	fmt.Println(v)
	//	fmt.Println(i)
	//}
	fmt.Println(x)
}

func findAnagrams(words []string) map[string][]string {
	words = toLowerCase(words)

	res := make(map[string][]string, len(words))
	for _, word := range words {
		sortStr := sortString(word)
		res[sortStr] = append(res[sortStr], word)
	}
	res = firstWord(res)

	return res
}

func firstWord(words map[string][]string) map[string][]string {
	res := make(map[string][]string, len(words))
	for _, word := range words {
		if len(word) == 1 {
			continue
		}
		tmp := word[0]
		sort.Strings(word)
		res[tmp] = append(res[tmp], word...)
	}
	return res
}

func toLowerCase(words []string) []string {
	wordMap := make(map[string]struct{}, len(words))
	var uniqueWords []string
	for _, word := range words {
		word = strings.ToLower(word)
		if _, ok := wordMap[word]; !ok {
			wordMap[word] = struct{}{}
			uniqueWords = append(uniqueWords, word)
		}
	}
	return uniqueWords
}

func sortString(s string) string {
	runes := []rune(s)
	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})
	return string(runes)
}
