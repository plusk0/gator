package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/plusk0/gator/internal/database"
)

func handlerDatabase(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]
	currentUser := database.CreateUserParams{uuid.New(), time.Now(), time.Now(), name}
	user, err := s.db.CreateUser(context.Background(), currentUser)
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}
	comm := command{"login", []string{name}}
	handlerLogin(s, comm)

	fmt.Printf("User %v created successfully!", user)
	return nil
}
