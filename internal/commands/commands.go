package commands

import (
	"fmt"

	"github.com/EyuAtske/Agrregator/internal/config"
)

type State struct{
	Cofg *config.Config
}

type Command struct{
	Name string
	Args []string
}

type Commands struct{
	Handlers map[string]func(*State, Command) error
}

func HandlerLogin(s *State, cmd Command) error{
	if len(cmd.Args) < 1 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}
	err := s.Cofg.SetUser(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("unable to set user: %w", err)
	}
	fmt.Println("The user name has been set")
	return nil
}

func (c *Commands) Run(s *State, cmd Command) error{
	handler, ok := c.Handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("command not found: %s", cmd.Name)
	}
	return handler(s, cmd)
}

func (c *Commands) Register(name string, f func(*State, Command) error){
	c.Handlers[name] = f
}