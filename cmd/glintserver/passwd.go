package main

import (
	"errors"
	"fmt"

	"github.com/urfave/cli"
)

func cliPasswd(c *cli.Context) error {
	var err error
	_, storage, err := setup(c, false)
	if err != nil {
		return err
	}
	defer cleanup(nil, storage)
	user := c.String("user")
	if user == "" {
		return errors.New("User not specified")
	}
	password, err := inputPassword("Enter new password: ", false)
	if err != nil {
		return errors.New("Error inputting password")
	}
	err = storage.ChangePassword(user, password)
	if err != nil {
		return err
	}
	fmt.Printf("Password updated for user '%s'\n", user)
	return nil
}
