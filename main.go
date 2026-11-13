package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/queueunit/Bloggo/internal/config"
	"github.com/queueunit/Bloggo/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	config *config.Config
	db     *database.Queries
}

func main() {

	// Start by getting the config from .gatorconfig.json
	// Create empty commands struct and populate with the register method
	x := config.Read()
	dbURL := x.Db_url
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Could not open database at: %v", dbURL)
	}

	newDB := database.New(db)

	c := commands{available: make(map[string]func(s *state, cmd command) error)}
	myState := &state{config: &x, db: newDB}
	c.register("login", handlerLogin)
	c.register("register", register)
	c.register("reset", reset)
	c.register("users", getUsers)

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
