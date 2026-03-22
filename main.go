package main

import (
	"fmt"
	"os"

	"github.com/EyuAtske/Agrregator/internal/config"
)

type state struct{
	cofg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Unable to read config")
		return
	}
	st := state{cofg: cfg}
	cmds := commands{
		Handlers: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	args := os.Args
	if len(args) == 1 {
		fmt.Println("not enough arguments were provided")
		os.Exit(1)
	}
	if len(args) == 2 {
		fmt.Println("a username is required")
		os.Exit(1)
	}
	cmd := command{Name: args[1], Args: args[2:]}
	err = cmds.run(&st, cmd)
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		os.Exit(1)
	}
}
