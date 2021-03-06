package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func mainInternal() int {
	if len(os.Args) < 3 {
		fmt.Println("Usage: envfile foo.env command")
		return 1
	}

	fp, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(fmt.Sprintf("envfile: %s", err.Error()))
		return 1
	}

	defer fp.Close()

	reader := bufio.NewReaderSize(fp, 4096)
	for {
		line, _, err := reader.ReadLine()
		s := strings.SplitN(string(line), "=", 2)
		if len(s) == 2 {
			os.Setenv(s[0], s[1])
		} else if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(fmt.Sprintf("envfile: %s", err.Error()))
			return 1
		}
	}

	path, err := exec.LookPath(os.Args[2])
	if err != nil {
		fmt.Println(fmt.Sprintf("envfile: %s", err.Error()))
		return 1
	}

	err = syscall.Exec(path, os.Args[2:], os.Environ())
	if err != nil {
		fmt.Println(fmt.Sprintf("envfile: %s", err.Error()))
		return 1
	}

	return 0
}

func main() {
	os.Exit(mainInternal())
}
