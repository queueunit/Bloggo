package main

import (
	"errors"
	"fmt"
)

type command struct {
	name      string
	arguments []string
}

type commands struct {
	available map[string]func(s *state, cmd command) error
}

func (c *commands) run(s *state, cmd command) error {

	if fxn, exists := c.available[cmd.name]; exists {
		err := fxn(s, cmd)
		if err != nil {
			return err
		}
	}
	return errors.New("")
}

func (c *commands) register(name string, f func(*state, command) error) {

	if _, exists := c.available[name]; exists {
		fmt.Printf("Function name %v exists, please choose a different name\n", name)
		return
	} else {
		c.available[name] = f
	}
}
