package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := handleEnvs(env)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	command.Env = os.Environ()

	err = command.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		return 1
	}

	return 0
}

func handleEnvs(env Environment) error {
	for envName, envItem := range env {
		if envItem.NeedRemove {
			err := os.Unsetenv(envName)
			if err != nil {
				return err
			}
		}

		err := os.Setenv(envName, envItem.Value)
		if err != nil {
			return err
		}
	}

	return nil
}
