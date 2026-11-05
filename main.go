package main

import (
	"log"
	"os"

	"github.com/queueunit/Bloggo/internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	// Start by getting the config from .gatorconfig.json
	// Create empty commands struct and populate with the register method
	x := config.Read()
	c := commands{available: make(map[string]func(s *state, cmd command) error)}
	myState := &state{config: &x}
	c.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("Error, missing command line argument")
	}

	cliCommand := os.Args[1]
	cliArguments := os.Args[2:]

	commandStruct := command{
		name:      cliCommand,
		arguments: cliArguments,
	}

	c.run(myState, commandStruct)

	return
}
