package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nassibnassar/goconfig/ini"
	"golang.org/x/crypto/ssh/terminal"
)

// inputPassword gets keyboard input from the user with terminal echo disabled.
// This function is intended for inputting passwords.  It prints a specified
// prompt before the input, and can optionally input the password a second time
// for confirmation.  The password is returned, or an error if there was a
// problem reading the input or (in the case of a confirmation input) if the
// two passwords did not match.  SIGINT is disabled during the input, to avoid
// leaving the terminal in a no-echo state.
func inputPassword(prompt string, confirm bool) (string, error) {
	// Ignore SIGINT, to avoid leaving terminal in no-echo state.
	signal.Ignore(os.Interrupt)
	defer signal.Reset(os.Interrupt)
	// Read the input.
	fmt.Print(prompt)
	p, err := terminal.ReadPassword(syscall.Stdin)
	fmt.Println("")
	if err != nil {
		return "", err
	}
	// Read the input again to confirm.
	if confirm {
		fmt.Print("(Confirming) " + prompt)
		q, err := terminal.ReadPassword(syscall.Stdin)
		fmt.Println("")
		if err != nil {
			return "", err
		}
		if string(p) != string(q) {
			return "", errors.New("Passwords do not match")
		}
	}
	// Return password.
	return string(p), nil
}

// readConfig reads the specified configuration file and returns its contents
// as a Config, or exits with an error message.
func oldReadConfig(file string) *ini.Config {
	var glintconfig *ini.Config
	var err error
	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		glintconfig = ini.NewConfig()
	} else {
		glintconfig, err = ini.NewConfigFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, os.Args[0]+
				": error reading configuration file '%s': %v\n",
				file, err)
			os.Exit(1)
		}
	}
	return glintconfig
}
