package main

import (
	"context"
	"fmt"
	"time"
	"os"

	"github.com/EyuAtske/Agrregator/internal/database"
	"github.com/google/uuid"
)

type command struct{
	Name string
	Args []string
}

type commands struct{
	Handlers map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error{
	if len(cmd.Args) < 1 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}
	_, err :=s.db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		fmt.Println("There is no user with the provided name")
		os.Exit(1)
	}
	err = s.cfg.SetUser(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("unable to set user: %w", err)
	}
	fmt.Println("The user name has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error{
	if len(cmd.Args) < 1 {
		return fmt.Errorf("the register handler expects a single argument, the username")
	}
	_, err :=s.db.GetUser(context.Background(), cmd.Args[0])
	if err == nil {
		fmt.Print(err)
		os.Exit(1)
	}
	err = s.cfg.SetUser(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("unable to set user: %w", err)
	}
	arg := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.Args[0],
	}
	_, err = s.db.CreateUser(context.Background(), arg)
	if err != nil {
		return fmt.Errorf("unable to create user: %w", err)
	}
	fmt.Println("User was created")
	return nil
}
func (c *commands) run(s *state, cmd command) error{
	handler, ok := c.Handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("command not found: %s", cmd.Name)
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error){
	c.Handlers[name] = f
}