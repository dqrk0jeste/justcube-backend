package main

import (
	"database/sql"
	"log"

	"github.com/dqrk0jeste/letscube-backend/api"
	database "github.com/dqrk0jeste/letscube-backend/database/sqlc"
	"github.com/dqrk0jeste/letscube-backend/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	connection, err := sql.Open(config.DatabaseDriver, config.DatabaseSource)
	if err != nil {
		log.Fatal(err)
	}

	database := database.New(connection)
	server, err := api.CreateServer(config, database)
	if err != nil {
		log.Fatal("error making a server: ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("error while starting a server", err)
	}
}
