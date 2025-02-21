package main

import (
	"testing"
	"fmt"
	"os"
	"bytes"
)

func TestCommandHelp(t *testing.T) {
	cliCommands = map[string]cliCommand{
		"exit": {name: "exit", description: "Exit the Pokedex", callback: commandExit},
		"help": {name: "help", description: "Displays a help message", callback: commandHelp},
	}

	var output bytes.Buffer
	fmt.Fprintln(&output, "Welcome to the Pokedex!\nUsage:")
	for _, cmd := range cliCommands {
		fmt.Fprintf(&output, "%s: %s\n", cmd.name, cmd.description)
	}

	expected := output.String()

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	commandHelp()

	w.Close()
	var testOutput bytes.Buffer
	testOutput.ReadFrom(r)
	os.Stdout = oldStdout // Restore stdout

	if testOutput.String() != expected {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expected, testOutput.String())
	}
}

func TestCommandExit(t *testing.T) {
	exited := false

	exitFunc = func(code int) {
		exited = true
	}

	commandExit()

	if !exited {
		t.Errorf("Expected commandExit to call exitFunc, but it did not")
	}

	exitFunc = os.Exit // Restore original exit function
}

