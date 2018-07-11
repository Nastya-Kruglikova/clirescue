package cmdutil

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ReadLine reads user input from command line.
func ReadLine() string {
	buf := bufio.NewReader(os.Stdin)
	line, err := buf.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	return strings.TrimSpace(string(line))
}

// Silence makes the shell not echo typed text back in terminal.
// Useful, when asking a user to type sensitive stuff like credentials.
func Silence() {
	runCommand(exec.Command("stty", "-echo"))
}

// Unsilence makes the shell echo typed text in the terminal.
// Should be run after Silence()
func Unsilence() {
	runCommand(exec.Command("stty", "echo"))
}

// runCommand takes exec.Command() as a parameter and executes it.
func runCommand(command *exec.Cmd) {
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Run()
}
