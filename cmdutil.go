package main

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
)

func readLine() (string, error) {
	r := bufio.NewReader(os.Stdin)
	line, err := r.ReadString('\n')
	return strings.TrimSpace(line), err
}

func silence() {
	runCommand(exec.Command("stty", "-echo"))
}

func unsilence() {
	runCommand(exec.Command("stty", "echo"))
}

func runCommand(command *exec.Cmd) {
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Run()
}
