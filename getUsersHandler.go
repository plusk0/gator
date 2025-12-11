package main

import (
	"context"
	"fmt"
)

func getUsers(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	current := s.cfg.CurrentUserName

	for _, v := range users {
		if current == v.Name {
			fmt.Println(v.Name, "(current)")
		} else {
			fmt.Println(v.Name)
		}
	}

	return nil
}
