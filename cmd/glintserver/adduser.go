package main

import (
	"errors"
	"fmt"

	"github.com/urfave/cli"
)

func cliAddUser(c *cli.Context) error {
	_, storage, err := setup(c, false)
	if err != nil {
		return err
	}
	defer cleanup(nil, storage)
	// Get user name.
	user := c.String("user")
	if user == "" {
		return errors.New("User not specified")
	}
	// Get full name.
	fullname := c.String("fullname")
	// Get email address.
	email := c.String("email")
	// Get password.
	password, err := inputPassword("Enter new password: ", false)
	if err != nil {
		return errors.New("Error inputting password")
	}
	// Add to Postgres.
	err = storage.AddPerson(user, fullname, email, password)
	if err != nil {
		return err
	}
	fmt.Printf("User '%s' added\n", user)
	return nil
}
