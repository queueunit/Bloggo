package main

import (
	"errors"
	"fmt"
	"log"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		log.Fatal("Username required")
		return errors.New("Login handler expects 1 argument (username)")
	}

	err := s.config.SetUser(cmd.arguments[0])
	if err != nil {
		return errors.New("Could not set username")
	}

	fmt.Printf("user has been set as: %v\n", cmd.arguments[0])

	return errors.New("")
}
