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

func reset(s *state, cmd command) error {

	contextBkg := context.Background()
	err := s.db.Reset(contextBkg)
	if err != nil {
		log.Fatal("Could not reset table users")
	}

	fmt.Println("Table 'users' reset")
	return errors.New("")
}

func getUsers(s *state, cmd command) error {

	contextBkg := context.Background()

	users, err := s.db.GetUsers(contextBkg)
	if err != nil {
		log.Fatal("Could not get column name from user table")
	}

	for _, x := range users {
		if x == s.config.CurrentUser {
			fmt.Printf("* %v (current)\n", x)
		} else {
			fmt.Printf("* %v\n", x)
		}
	}

	return errors.New("")
}

//func addFeed(s *state, cmd command) error {
//
//	current := s.config.CurrentUser
//	ctx := context.Background()
//	user, err := s.db.GetUser(ctx, current)
//	if err != nil {
//		log.Fatal("Could not fetch info for user: %v", current)
//	}
//
//	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
//		ID:        uuid.New(),
//		CreatedAt: time.Now(),
//		UpdatedAt: time.Now(),
//		Name:      cmd.arguments[0],
//		Url:       cmd.arguments[1],
//		UserID:    user.ID,
//	})
//	if err != nil {
//		log.Fatal("Could not add feed %v", err)
//	}
//
//	// Create a feed follow immediately after the user creates a feed
//	var feedFollow = database.CreateFeedFollowParams{
//		ID:        uuid.New(),
//		CreatedAt: time.Now(),
//		UpdatedAt: time.Now(),
//		UserID:    feed.UserID,
//		FeedID:    feed.ID,
//	}
//
//	_, err = s.db.CreateFeedFollow(ctx, feedFollow)
//	if err != nil {
//		log.Fatal("Could not add corresponding feedFollow")
//	}
//
//	fmt.Println(feed.Name)
//	fmt.Println(feed.Url)
//	return errors.New("")
//}

func addFeed(s *state, cmd command, user database.User) error {

	ctx := context.Background()

	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
		Url:       cmd.arguments[1],
		UserID:    user.ID,
	})
	if err != nil {
		log.Fatal("Could not add feed %v", err)
	}

	// Create a feed follow immediately after the user creates a feed
	var feedFollow = database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    feed.UserID,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(ctx, feedFollow)
	if err != nil {
		log.Fatal("Could not add corresponding feedFollow")
	}

	fmt.Println(feed.Name)
	fmt.Println(feed.Url)
	return errors.New("")
}
