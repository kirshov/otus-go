package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type envItem struct {
	file  string
	value string
}

func TestReadDir(t *testing.T) {
	i := []envItem{
		{file: "BAR", value: "bar"},
		{file: "EMPTY", value: ""},
		{file: "FOO", value: "   foo\nwith new line"},
		{file: "HELLO", value: "\"hello\""},
		{file: "UNSET", value: ""},
	}

	env, err := ReadDir("./testdata/env")

	require.NoError(t, err)

	for _, v := range i {
		require.Equal(t, v.value, env[v.file].Value)
	}
}
