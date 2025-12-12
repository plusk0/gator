package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/plusk0/gator/internal/config"
	"github.com/plusk0/gator/internal/database"
)

type state struct {
	db            *database.Queries
	currentOffset int
	cfg           *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("error reading config:", err)
	}

	dbs, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatal("Failed to open db:", err)
	}

	dbQueries := database.New(dbs)
	programState := &state{
		cfg: &cfg,
	}
	programState.db = dbQueries
	programState.currentOffset = 0

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerDatabase)
	cmds.register("reset", handlerReset)
	cmds.register("users", getUsers)
	cmds.register("agg", aggHandler)
	cmds.register("feeds", feedListHandler)

	cmds.register("addfeed", middlewareLoggedIn(addFeedHandler))
	cmds.register("follow", middlewareLoggedIn(followHandler))
	cmds.register("following", middlewareLoggedIn(followingHandler))
	cmds.register("unfollow", middlewareLoggedIn(unfollowHandler))
	cmds.register("browse", middlewareLoggedIn(browseHandler))
	cmds.register("next", middlewareLoggedIn(nextHandler))

	cmds.register("exit", func(s *state, cmd command) error { os.Exit(0); return nil })

	// Start REPL
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break // Exit on EOF (Ctrl+D)
		}
		input := scanner.Text()
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			continue
		}

		// Parse input into command and args
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}
		cmdName := parts[0]
		cmdArgs := parts[1:]

		// Run the command
		err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
		}
	}
}
