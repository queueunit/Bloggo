package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/queueunit/Bloggo/internal/database"
)

//func addFeedFollow(s *state, cmd command) error {
//
//	if len(cmd.arguments) == 0 {
//		log.Fatal("Need URL argument")
//	}
//
//	currentUser := s.config.CurrentUser
//	ctx := context.Background()
//	feed, err := s.db.LookUpFeedByURL(ctx, cmd.arguments[0])
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	user, err := s.db.GetUser(ctx, currentUser)
//	if err != nil {
//		log.Fatal("Could not fetch info for user: %v", currentUser)
//	}
//
//	feedFollow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
//		ID:        uuid.New(),
//		CreatedAt: time.Now(),
//		UpdatedAt: time.Now(),
//		UserID:    user.ID,
//		FeedID:    feed.ID,
//	})
//	if err != nil {
//		log.Fatal("Could not add new Feed Follow")
//	}
//
//	fmt.Println("User: %v now following %v", user.Name, feedFollow.FeedName)
//	return nil
//
//}

func addFeedFollow(s *state, cmd command, user database.User) error {

	if len(cmd.arguments) == 0 {
		log.Fatal("Need URL argument")
	}

	ctx := context.Background()
	feed, err := s.db.LookUpFeedByURL(ctx, cmd.arguments[0])
	if err != nil {
		log.Fatal(err)
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

//func getFollowsForUser(s *state, cmd command) error {
//
//	ctx := context.Background()
//	user, err := s.db.GetUser(ctx, s.config.CurrentUser)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	userFollows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
//	if err != nil {
//		log.Fatal("Could not fetch follows for current user")
//	}
//
//	fmt.Printf("User %v Following:\n", user.Name)
//	for _, f := range userFollows {
//		fmt.Println(f.Blog)
//	}
//
//	return nil
//}

func getFollowsForUser(s *state, cmd command, user database.User) error {

	ctx := context.Background()

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

func deleteFeedFollow(s *state, cmd command, user database.User) error {

	if len(cmd.arguments) == 0 {
		log.Fatal("unfollow command requires one argument (URL)")
	}

	ctx := context.Background()

	feed, err := s.db.LookUpFeedByURL(ctx, cmd.arguments[0])
	if err != nil {
		log.Fatal("Could not fetch feed information")
	}

	fmt.Println("User ID:   %v", user.ID)
	fmt.Println("Feed ID:   %v", feed.ID)
	err = s.db.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{UserID: user.ID, UserID_2: feed.ID})

	return nil
}
