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
			if _, err := strconv.Atoi(char); err == nil {
				return "", ErrInvalidString
			}
		}

		var resultIteration strings.Builder
		resultIteration.WriteString(char)

		if i < length-1 {
			next := string(sArr[i+1])
			countInt, err := strconv.Atoi(next)

			if err == nil {
				i++
				if countInt == 0 {
					continue
				}

				if countInt > 0 {
					resultIteration.WriteString(strings.Repeat(char, countInt-1))
				}
			}
		}

		result.WriteString(resultIteration.String())
	}

	return result.String(), nil
}

func getChar(a []rune, i int) (string, bool) {
	isSlashed := false
	char := string(a[i])

	if char == "\\" {
		l := len(a)
		if l > i+1 {
			isSlashed = true
			char = string(a[i+1])
		}
	}

	return char, isSlashed
}
