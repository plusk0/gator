package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/plusk0/gator/internal/config"
	"github.com/plusk0/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
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

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
