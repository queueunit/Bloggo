package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"time"

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

func scrapeFeeds(s *state, cmd command) error {

	ctx := context.Background()
	next, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		log.Fatal("Could not fetch feed to update")
	}

	err = s.db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		ID:        next.ID,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Fatal("Could not mark feed as updated")
	}

	feed, err := fetchFeed(ctx, next.Url)
	if err != nil {
		log.Fatal("Could not fetch requested feed from %v", next.Url)
	}

	for _, x := range feed.Channel.Item {
		fmt.Println(html.UnescapeString(x.Title))
	}

	return nil
}
