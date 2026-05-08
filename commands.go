package main

import (
	"errors"
	"fmt"

	"github.com/hussmaster/gogetter/internal/config"
)

// State struct that holds a pointer to a config Struct
type state struct {
	configFile *config.Config
}

// Command struct that contains a name as a string type and args as a slice of strings
type command struct {
	name string
	args []string
}

// Login command function
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("no arguments")
	} else if len(cmd.args) > 1 {
		return errors.New("too many arguments")
	}
	//Get username
	userName := cmd.args[0]
	//Set username to json file
	err := s.configFile.SetUser(userName)
	if err != nil {
		return fmt.Errorf("could not set username: %w", err)
	}
	fmt.Printf("Username: %s has been set\n", userName)
	return nil
}

type commands struct {
	registry map[string]func(*state, command) error
}

// Function to run a command if it exists in the command registry
func (c *commands) run(s *state, cmd command) error {
	//Checks registry for command name
	value, ok := c.registry[cmd.name]
	if ok {
		//Runs command
		return value(s, cmd)
	}
	return fmt.Errorf("command: %v not in command registry", cmd.name)
}

// Function to register a command in the command registry
func (c *commands) register(name string, f func(*state, command) error) {
	//Registers the name of the command by the function name
	c.registry[name] = f
}
