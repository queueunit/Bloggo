package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/queueunit/Bloggo/internal/database"
)

func getUsers_DELETE_THIS(s *state, cmd command) error {

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

func addFeed_DELETE_THIS(s *state, cmd command) error {

	current := s.config.CurrentUser
	ctx := context.Background()
	user, err := s.db.GetUser(ctx, current)
	if err != nil {
		log.Fatal("Could not fetch info for user: %v", current)
	}

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

	fmt.Println(feed.Name)
	fmt.Println(feed.Url)
	return errors.New("")
}

func addFeedFollow(s *state, cmd command) error {

	if len(cmd.arguments) == 0 {
		log.Fatal("Need URL argument")
	}

	currentUser := s.config.CurrentUser
	ctx := context.Background()
	feed, err := s.db.LookUpFeedByURL(ctx, cmd.arguments[0])
	if err != nil {
		log.Fatal(err)
	}

	user, err := s.db.GetUser(ctx, currentUser)
	if err != nil {
		log.Fatal("Could not fetch info for user: %v", currentUser)
	}

	feedFollow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		log.Fatal("Could not add new Feed Follow")
	}

	fmt.Println("User: %v now following %v", user.Name, feedFollow.FeedName)
	return nil

}

func getFollowsForUser(s *state, cmd command) error {

	ctx := context.Background()
	user, err := s.db.GetUser(ctx, s.config.CurrentUser)
	if err != nil {
		log.Fatal(err)
	}

	userFollows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		log.Fatal("Could not fetch follows for current user")
	}

	fmt.Printf("User %v Following:\n", user.Name)
	for _, f := range userFollows {
		fmt.Println(f.Blog)
	}

	return nil
}
