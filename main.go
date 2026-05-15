package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/hussmaster/gogetter/internal/config"
	"github.com/hussmaster/gogetter/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading confg: %v\n", err)
	}
	//fmt.Printf("Read config: %+v\n", cfg)

	//Attempt db connection
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error opening database connection to: %v\n", err)
	}
	//Create new database connection
	dbQueries := database.New(db)

	//State struct
	s := &state{
		configFile: &cfg,
		db:         dbQueries,
	}

	//Create empty struct of commands
	//Initialize map inside of struct
	commandList := &commands{
		registry: make(map[string]func(*state, command) error),
	}

	//Register commands
	commandList.register("login", handlerLogin)
	commandList.register("register", handlerRegister)
	commandList.register("reset", handlerReset)
	commandList.register("users", handlerUsers)
	commandList.register("agg", handlerAgg)
	commandList.register("addfeed", handlerAddFeed)
	commandList.register("feeds", handlerGetFeeds)

	//Makes sure os.Args is greater than 2
	if len(os.Args) < 2 {
		log.Fatal("not enough arguments\n")
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
