package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/queueunit/Bloggo/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		log.Fatal("Username required")
		return errors.New("Login handler expects 1 argument (username)")
	}

	contextBkg := context.Background()
	_, err := s.db.GetUser(contextBkg, cmd.arguments[0])
	if err != nil {
		fmt.Println("User does not exist")
		os.Exit(1)
	}

	err = s.config.SetUser(cmd.arguments[0])
	if err != nil {
		return errors.New("Could not set username")
	}

	fmt.Printf("user has been set as: %v\n", cmd.arguments[0])

	return errors.New("")
}

func register(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		log.Fatal("Username required")
	}

	contextBkg := context.Background()

	// Check to see if user exists before SQL query
	_, err := s.db.GetUser(contextBkg, cmd.arguments[0])
	if err == nil {
		fmt.Println("User already exists")
		os.Exit(1)
	}

	user, err := s.db.CreateUser(contextBkg, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
	})
	if err != nil {
		log.Fatal("Could not create user with name: %v", cmd.arguments[0])
	}

	err = s.config.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("failed to set current user: %v", err)
	}
	fmt.Print(`
	ID: %v
	CreatedAt: %v
	UpdatedAt: %v
	Name: %v`, user.ID, user.CreatedAt, user.UpdatedAt, user.Name)

	return errors.New("")
}
