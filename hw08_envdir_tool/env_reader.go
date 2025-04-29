package main

import (
	"bufio"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	readDir, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	result := make(Environment)

	for _, file := range readDir {
		if !file.IsDir() {
			v, err := parseFile(dir + "/" + file.Name())
			if err != nil {
				continue
			}

			result[file.Name()] = v
		}
	}

	// Place your code here
	return result, nil
}

func parseFile(path string) (EnvValue, error) {
	file, err := os.Open(path)
	if err != nil {
		return EnvValue{}, err
	}

	defer file.Close()

	s := bufio.NewScanner(file)
	s.Scan()

	value := s.Text()
	value = strings.TrimRight(value, " 	\n")
	value = strings.ReplaceAll(value, "\x00", "\n")

	return EnvValue{value, value != ""}, nil
}
