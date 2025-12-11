package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/plusk0/gator/internal/database"
)

func addFeedHandler(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s [Site Name] [URL]", cmd.Name)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]
	userID := user.ID

	feed := database.AddFeedParams{
		uuid.New(),
		time.Now(),
		time.Now(),
		name,
		url,
		userID,
	}

	_, err := s.db.AddFeed(context.Background(), feed)
	if err != nil {
		return err
	}

	fmt.Printf("Feed source %v at %v created successfully!", name, url)

	wrapped := middlewareLoggedIn(followHandler)
	err = wrapped(s, command{Name: "follow", Args: []string{url}})
	if err != nil {
		return err
	}
	return nil
}

func feedListHandler(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, v := range feeds {
		fmt.Println("Name:", v.Name)
		fmt.Println("Url:", v.Url)
		user, err := s.db.GetUserID(context.Background(), v.UserID)
		if err != nil {
			return err
		}
		fmt.Println("Added By:", user.Name)

	}

	return nil
}
