Basic guided project for implementing a web aggregator with Postresql in Go (using Goose, sqlc)
Dependencies:
postgres, go

Binary can be built from the directory via /go build and installed via /go install

The program expects a config file in the user home dir ~/.gatorconfig.json which defines the postgres-interface
Example:
{"db_url":"postgres://postgres:@localhost:5432/gator?sslmode=disable","current_user_name":"new"}
The username can be changed later via the command register [name]

Basic commands:
login "username"
register "username"
reset (Warning- this resets all aggregated posts/feeds/users (mainly for testing))

Info-commands:
users (shows all user information)
feeds (same for feeds)
following (shows feeds followed by current user, has to be logged in)

Practical commands:
addfeed "name" "url" (adds the url as a new RSS source, automatically follows)
agg "time" (starts web scraper, aggregates every X time. Example Format: 1h15m30s)
follow/unfollow "url" (allows multiple users to follow feeds of others)
browse "number" (shows x amount of recent posts from followed sources)


To-do for practice:

    Add sorting and filtering options to the browse command
    Add pagination to the browse command
    Add concurrency to the agg command so that it can fetch more frequently
    Add a search command that allows for fuzzy searching of posts
    Add bookmarking or liking posts
    Add a TUI that allows you to select a post in the terminal and view it in a more readable format (either in the terminal or open in a browser)
    Add an HTTP API (and authentication/authorization) that allows other users to interact with the service remotely
    Write a service manager that keeps the agg command running in the background and restarts it if it crashes



Basic stuff to remind myself:
Startup postgresql (on archlinux):

sudo systemctl start postgresql
sudo -u postgres psql
\c gator

Database Migrations with goose:
cd sql/schema
goose postgres "connection string" up/down

sqlc for secure go/sql code generation
sqlc generate (from project dir)
