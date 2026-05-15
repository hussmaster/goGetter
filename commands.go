package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hussmaster/gogetter/internal/config"
	"github.com/hussmaster/gogetter/internal/database"
)

// State struct that holds a pointer to a config Struct
// also holds pointer to database query connection
type state struct {
	configFile *config.Config
	db         *database.Queries
}

// Command struct that contains a name as a string type and args as a slice of strings
type command struct {
	name string
	args []string
}

// Login command function
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("no arguments\n")
	} else if len(cmd.args) > 1 {
		return errors.New("too many arguments\n")
	}
	//Get username
	userName := cmd.args[0]
	//Set username to json file
	err := s.configFile.SetUser(userName)
	if err != nil {
		return fmt.Errorf("could not set username: %w\n", err)
	}
	fmt.Printf("Username: %s has been set\n", userName)

	//Check if user exists in the database before logging in
	_, err = s.db.GetUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("unable to login as: %v doesn't exist in the database\n", err)
	}
	return nil
}

// Registers user in database
func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("no arguments\n")
	} else if len(cmd.args) > 1 {
		return errors.New("too many arguments\n")
	}

	//Get username from cmdline arguments
	userName := cmd.args[0]
	//Set current user in the config file
	err := s.configFile.SetUser(userName)
	if err != nil {
		return errors.New("unable to set current username in the config file\n")
	}
	//Check if username already exists in the database
	_, err = s.db.GetUser(context.Background(), userName)
	if err == nil {
		log.Fatalf("username: %v already exists\n", err)
	}
	//Register username in database
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      userName,
	})
	if err != nil {
		return err
	}
	fmt.Printf("username: %v was created in the database\n", userName)
	fmt.Printf("%+v\n", user)

	return nil
}

// Resets state of the database
func handlerReset(s *state, cmd command) error {
	err := s.db.DelUsers(context.Background())
	if err != nil {
		return errors.New("database was unable to be reset")
	}
	return nil
}

// Display all users in database, note which one is currently logged in
func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return errors.New("unable to establish connection with database")
	}
	curUser := s.configFile.CurrentUserName
	for _, user := range users {
		if user.Name == curUser {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}
	return nil
}

// Aggregates returned xml into database
func handlerAgg(s *state, cmd command) error {
	tempURL := "https://www.wagslane.dev/index.xml"
	xmlRSS, err := fetchFeed(context.Background(), tempURL)
	if err != nil {
		return fmt.Errorf("failed to fetch url with error: %w", err)
	}
	fmt.Printf("%+v\n", xmlRSS)
	return nil
}

// Adds feed to database
func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return errors.New("not enough arguments\n")
	}
	//Get current user from the config file
	curUser := s.configFile.CurrentUserName
	//Get username from database using the username from the configfile
	curDBUser, err := s.db.GetUser(context.Background(), curUser)
	if err != nil {
		return fmt.Errorf("error quering user: %w", err)
	}
	//Create feed
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    curDBUser.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating the feed with error: %w", err)
	}
	fmt.Printf("feed for %v created in the database\n", cmd.args[0])
	//Return created feed record
	fmt.Printf("%+v\n", feed)
	return nil
}

// Lists feeds in database
func handlerGetFeeds(s *state, cmd command) error {
	if len(cmd.args) > 1 {
		return errors.New("too many arguments. feeds takes 0 cmdline arguments\n")
	}
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("unable to query feeds table with error: %w", err)
	}
	for _, feed := range feeds {
		feedName := feed.Name
		feedURL := feed.Url
		feedUserUserID := feed.UserID
		//Query users db with uuid in feeds table
		user, err := s.db.GetFeedUser(context.Background(), feedUserUserID)
		if err != nil {
			return fmt.Errorf("unable to query users table with error: %w", err)
		}
		feedUser := user.Name
		//fmt.Printf("Feed name: %v Feed URL: %v Feed created by: %v\n", feedName, feedURL, feedUser)
		fmt.Printf("%v\n%v\n%v\n", feedName, feedURL, feedUser)
	}
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
