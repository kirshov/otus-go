package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	sArr := []rune(s)
	length := len(sArr)
	var result strings.Builder

	for i := 0; i < length; i++ {
		char, isSlashed := getChar(sArr, i)
		if isSlashed {
			i++
		} else {
			// сдвиг всегда на 2 символа минимум, если пришла неэкранированная цифра - ошибка
			if _, err := strconv.Atoi(char); err == nil {
				return "", ErrInvalidString
			}
		}

		count, inc := getRepeatCount(sArr, i)
		i += inc
		if count < 0 {
			continue
		}

		result.WriteString(char)
		if count > 0 {
			result.WriteString(strings.Repeat(char, count))
		}
	}

	return result.String(), nil
}

// Вернуть текущий символ, учитывая экранирование (если экранирование есть, то следующий за ним).
func getChar(a []rune, i int) (string, bool) {
	if a[i] == '\\' && i+1 < len(a) {
		return string(a[i+1]), true
	}

	return string(a[i]), false
}

// Вернуть количество повторений текущего символа и инкремент
func getRepeatCount(a []rune, i int) (count int, inc int) {
	if i+1 >= len(a) {
		return
	}

	next := string(a[i+1])
	countInt, err := strconv.Atoi(next)

	// если ошибки нет, значит следующий символ - число
	if err == nil {
		if countInt == 0 {
			inc = 1
			count = -1
		}

		if countInt > 0 {
			count = countInt - 1
			inc = 1
		}
	}

	return
}
