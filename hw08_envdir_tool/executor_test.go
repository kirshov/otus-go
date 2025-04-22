package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	envName := "FOO"
	envValue := "bar"
	envItem := EnvValue{
		Value: envValue,
	}
	env := make(Environment, 1)
	env[envName] = envItem

	code := RunCmd([]string{"printenv", envName}, env)

	require.Equal(t, envItem.Value, os.Getenv(envName))
	require.Equal(t, 0, code)
}
