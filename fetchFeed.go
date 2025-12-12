package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/plusk0/gator/internal/database"
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

func aggHandler(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s [timeInterval: 3h/2m/1s]", cmd.Name)
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Println("Collecting feeds every ", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err = webScrape(s)
		if err != nil {
			return err
		}
	}
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var data RSSFeed
	err = xml.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	data.Channel.Link = html.UnescapeString(data.Channel.Link)
	data.Channel.Description = html.UnescapeString(data.Channel.Description)

	for _, v := range data.Channel.Item {
		v.Link = html.UnescapeString(v.Link)
		v.Description = html.UnescapeString(v.Description)
	}
	return &data, nil
}

func webScrape(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	_, err = s.db.MarkFetchedFeeds(context.Background(), feed.ID)
	if err != nil {
		return err
	}
	data, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return nil // err
	}
	fmt.Println("Found data,...Logging")
	x := 1

	for _, v := range data.Channel.Item {
		publishedAt := sql.NullTime{}
		if x == 1 {
			fmt.Println(v)
			x = 2
		}
		if t, err := time.Parse(time.RFC1123Z, v.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}
		_, err = s.db.AddPost(context.Background(), database.AddPostParams{
			uuid.New(),
			time.Now().UTC(),
			time.Now().UTC(),
			v.Title,
			v.Link,
			sql.NullString{v.Description, true},
			publishedAt,
			feed.ID,
		})
	}
	return nil
}

func browseHandler(s *state, cmd command, user database.User) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %s [MaxPostNumber]", cmd.Name)
	}

	num := 1
	if len(cmd.Args) == 1 {
		var err error
		num, err = strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("usage: %s [MaxPostNumber as int]", cmd.Name)
		}
	}

	offset := s.currentOffset
	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(num),
		Offset: int32(offset),
	}

	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return err
	}

	for _, v := range posts {
		fmt.Printf("\n %v fetched %v : %v \n %v \n", v.FeedName, v.UpdatedAt.Format("2006-01-02 15:04"), v.Title, v.Description.String)
	}

	// Update the offset for the next page
	s.currentOffset += num
	fmt.Println(s.currentOffset)
	return nil
}

func nextHandler(s *state, cmd command, user database.User) error {
	if s.currentOffset == 0 {
		return fmt.Errorf("no previous browse command detected")
	}

	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(3),
		Offset: int32(s.currentOffset),
	}

	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return err
	}

	for _, v := range posts {
		fmt.Println(v.Title, v.Description)
	}

	s.currentOffset += 3

	return nil
}
