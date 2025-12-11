Basic guided project for implementing a web aggregator with Postresql in Go (using Goose, sqlc)
Dependencies:
postgres, go

Binary can be built from the directory via /go build and installed via /go install

The program expects a config file in the user home dir ~/.gatorconfig.json
Example:
{"db_url":"postgres://postgres:@localhost:5432/gator?sslmode=disable","current_user_name":"new"}
The username can be set later via the command register [name]



Basic stuff to remind myself:
Startup postgresql (on archlinux):

sudo systemctl start postgresql
sudo -u postgres psql
\c gator

Database Migrations with goose:
cd sql/schema
goose postgres "postgres://postgres:@localhost:5432/gator" up/down

sqlc for secure go/sql code generation
sqlc generate (from project dir)
