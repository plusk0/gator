package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/plusk0/gator/internal/database"
)

func followHandler(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <URL>", cmd.Name)
	}
	url := cmd.Args[0]

	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		return err
	}
	params := database.CreateFeedFollowParams{uuid.New(), time.Now(), time.Now(), user.ID, feed.ID}
	result, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Println(result, feed.Name, s.cfg.CurrentUserName)
	return nil
}

func followingHandler(s *state, cmd command, user database.User) error {
	name := user.Name
	output, err := s.db.GetFeedFollowsForUser(context.Background(), name)
	if err != nil {
		return err
	}
	fmt.Println(output)
	return nil
}

func unfollowHandler(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <URL>", cmd.Name)
	}
	url := cmd.Args[0]

	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		return err
	}
	params := database.UnFollowParams{user.ID, feed.ID}

	_, err = s.db.UnFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Println(feed.Name, s.cfg.CurrentUserName)
	return nil
}
