package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("args is empty")
		os.Exit(1)
	}

	envs, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Println("read dir error:", err)
		os.Exit(1)
	}

	os.Exit(RunCmd(os.Args[2:], envs))
}
