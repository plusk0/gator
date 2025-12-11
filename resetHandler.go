package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	err := s.db.ResetUser(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("Users reset successfully!")
	return nil
}
