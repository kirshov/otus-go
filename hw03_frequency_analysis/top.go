package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

const listSize = 10

var re = regexp.MustCompile(`[\p{L}-]+|-`)

type wordsCount struct {
	word  string
	count int
}

func Top10(s string) []string {
	words := re.FindAllString(strings.ToLower(s), -1)
	list := make(map[string]int)

	for _, word := range words {
		list[word]++
	}

	return getSortTopList(list)
}

func getSortTopList(list map[string]int) []string {
	sl := createWordsCountSlice(list)
	sl = sortList(sl)

	maxSize := listSize
	if len(sl) < listSize {
		maxSize = len(sl)
	}

	r := make([]string, 0, maxSize)
	for _, item := range sl[:maxSize] {
		r = append(r, item.word)
	}

	return r
}

func createWordsCountSlice(list map[string]int) []wordsCount {
	sl := make([]wordsCount, 0, len(list))

	for w, c := range list {
		if w == "-" {
			continue
		}

		sl = append(sl, wordsCount{w, c})
	}

	return sl
}

func sortList(sl []wordsCount) []wordsCount {
	sort.Slice(sl, func(i, j int) bool {
		if sl[i].count == sl[j].count {
			return sl[i].word < sl[j].word
		}

		return sl[i].count > sl[j].count
	})

	return sl
}
