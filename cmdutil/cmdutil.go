package cmdutil

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/GoBootcamp/clirescue/user"
)

var (
	inputFile   *os.File
	inputBuffer *bufio.Reader
	stdout      *os.File
)

func init() {
	inputFile = os.Stdin
	stdout = os.Stdout
}

// Credentials reads username and password and returns a pointer to a User structure with sered credentials
func Credentials() (*user.User, error) {
	fmt.Fprint(stdout, "Username: ")
	username, err := readLine()
	if err != nil {
		return nil, err
	}

	silence()

	fmt.Fprint(stdout, "Password: ")
	password, err := readLine()
	if err != nil {
		return nil, err
	}

	unsilence()

	currentUser := user.New()
	currentUser.SetLogin(username, password)
	return currentUser, nil
}

func readLine() (string, error) {
	buf := buffer()
	line, err := buf.ReadString('\n')
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(line), nil
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

func buffer() *bufio.Reader {
	if inputBuffer == nil {
		inputBuffer = bufio.NewReader(inputFile)
	}
	return inputBuffer
}
