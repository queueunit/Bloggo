package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/queueunit/Bloggo/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		log.Fatal("Could not create Get request to %v", feedURL)
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Unexpected status: %s\n", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body: %v\n", err)
	}

	var rss RSSFeed
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return &rss, nil

}

func agg(s *state, cmd command, time_between_reqs string) {

	interval, _ := time.ParseDuration(time_between_reqs)
	ticker := time.NewTicker(interval)

	fmt.Println("----- Collecting feeds every %v", interval)
	for ; ; <-ticker.C {
		scrapeFeeds(s, cmd)
	}

}

func showFeeds(s *state, cmd command) error {

	ctx := context.Background()

	feeds, err := s.db.ListFeeds(ctx)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Could not fetch feeds in database")
	}

	fmt.Println("---------- Feeds ----------")
	for _, f := range feeds {
		fmt.Println(f.Name)
		fmt.Println(f.Url)
		fmt.Println(f.Name_2)
	}

	return nil
}

//func scrapeFeeds(s *state, cmd command) error {
//
//	ctx := context.Background()
//	next, err := s.db.GetNextFeedToFetch(ctx)
//	if err != nil {
//		log.Fatal("Could not fetch feed to update")
//	}
//
//	err = s.db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
//		ID:        next.ID,
//		UpdatedAt: time.Now(),
//	})
//	if err != nil {
//		log.Fatal("Could not mark feed as updated")
//	}
//
//	feed, err := fetchFeed(ctx, next.Url)
//	if err != nil {
//		log.Fatal("Could not fetch requested feed from %v", next.Url)
//	}
//
//	for _, x := range feed.Channel.Item {
//		fmt.Println(html.UnescapeString(x.Title))
//	}
//
//	return nil
//}

func scrapeFeeds(s *state, cmd command) error {
	// get the next feed to fetch from the database
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	// mark it as fetched in database
	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID:        feed.ID,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	// fetch the rss feed from url
	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	for _, item := range rssFeed.Channel.Item {
		// parse published_at (2 formats)
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}
		if t, err := time.Parse(time.RFC1123, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		description := sql.NullString{
			String: item.Description,
			Valid:  item.Description != "",
		}

		_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       toNullString(item.Title),
			Url:         item.Link,
			Description: description,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			fmt.Printf("Error saving post: %v\n", err)
		}
	}

	return nil

}

func toNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "", // Change to 'true' if empty strings shouldn't be NULL
	}
}

func toNullTime(s string, layout string) (sql.NullTime, error) {
	if s == "" {
		return sql.NullTime{Valid: false}, nil
	}

	t, err := time.Parse(layout, s)
	if err != nil {
		return sql.NullTime{}, err
	}

	return sql.NullTime{
		Time:  t,
		Valid: true,
	}, nil
}
