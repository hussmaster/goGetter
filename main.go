package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hussmaster/gogetter/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading confg: %v", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	s := &state{
		configFile: &cfg,
	}

	//Create empty struct of commands
	//Initialize map inside of struct
	commandList := &commands{
		registry: make(map[string]func(*state, command) error),
	}

	//Register login command
	commandList.register("login", handlerLogin)

	//Makes sure os.Args is greater than 2
	if len(os.Args) < 2 {
		log.Fatalf("not enough arguments")
	}
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	//Creates command struct with command name and the arguments provided
	cmd := command{
		cmdName,
		cmdArgs,
	}
	//Runs command, checks for errors
	err = commandList.run(s, cmd)
	if err != nil {
		log.Fatal(err)
	}
}

//err = cfg.SetUser("ian")
//if err != nil {
//	log.Fatalf("couldn't set current user: %v", err)
//}

//cfg, err = config.Read()
//if err != nil {
//	log.Fatalf("error reading config: %v", err)
//}
//fmt.Printf("Read config again: %+v\n", cfg)
