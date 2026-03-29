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

func handlerReset(s *state, cmd command) error{
	err := s.db.DeleteUsers(context.Background())
	if err != nil{
		fmt.Println("Unable to reset users: %w", err)
		os.Exit(1)
	}
	return nil
}

func handlerUsers(s *state, cmd command) error{
	users, err := s.db.GetUsers(context.Background())
	if err != nil{
		fmt.Println("Unable to fetch users: %w", err)
		os.Exit(1)
	}
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("%s (current)\n", user.Name)
			continue
		}
		fmt.Println(user.Name)
	}
	
	return nil
}

func handlerFetch(s *state, cmd command) error{
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("unable to fetch feed: %w", err)
	}
	for _, item := range feed.Channel.Item {
		fmt.Printf("- %s\n", item.Title)
		fmt.Printf("  Link: %s\n", item.Link)
		fmt.Printf("  Description: %s\n", item.Description)
		fmt.Printf("  Published: %s\n", item.PubDate)
	}
	return nil
}

func handlerAddFeed(s *state, cmd command) error{
	if len(cmd.Args) < 2 {
		return fmt.Errorf("the addfeed handler expects two arguments, the feed name and the feed url")
	}
	name := cmd.Args[0]
	url := cmd.Args[1]
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("unable to fetch user: %w", err)
	}
	s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		user.ID,
		user.CreatedAt,
		user.UpdatedAt,
		name,
		url,
		user.ID,
	})
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