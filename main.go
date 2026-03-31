package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/EyuAtske/Agrregator/internal/config"
	"github.com/EyuAtske/Agrregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:datapost@localhost:5432/gator?sslmode=disable")
	if err != nil {
		fmt.Println("Unable to connect to database")
		return
	}
	defer db.Close()
	dbQueries := database.New(db)

	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Unable to read config")
		return
	}
	st := state{db: dbQueries, cfg: cfg}
	cmds := commands{
		Handlers: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	// for admins
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerFetch)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	args := os.Args
	cmd := command{}
	if len(args) == 1 {
		fmt.Println("not enough arguments were provided")
		os.Exit(1)
	}
	if len(args) == 2 {
		switch args[1] {
		case "reset":
			cmd = command{Name: args[1], Args: nil}
		case "users":
			cmd = command{Name: args[1], Args: nil}
		case "agg":
			cmd = command{Name: args[1], Args: nil}
		case "feeds":
			cmd = command{Name: args[1], Args: nil}
		case "following":
			cmd = command{Name: args[1], Args: nil}
		default:
			fmt.Println("a second argument is required")
			os.Exit(1)
		}
	}else{
		cmd = command{Name: args[1], Args: args[2:]}
	}
	err = cmds.run(&st, cmd)
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		os.Exit(1)
	}
}
