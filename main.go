package main

import (
	"fmt"
	"os"

	"github.com/EyuAtske/Agrregator/internal/commands"
	"github.com/EyuAtske/Agrregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Unable to read config")
		return
	}
	state := commands.State{Cofg: cfg}
	cmds := commands.Commands{
		Handlers: make(map[string]func(*commands.State, commands.Command) error),
	}
	cmds.Register("login", commands.HandlerLogin)
	args := os.Args
	if len(args) == 1 {
		fmt.Println("not enough arguments were provided")
		os.Exit(1)
	}
	if len(args) == 2 {
		fmt.Println("a username is required")
		os.Exit(1)
	}
	cmd := commands.Command{Name: args[1], Args: args[2:]}
	commands.HandlerLogin(&state, cmd)
	err = cmds.Run(&state, cmd)
}
