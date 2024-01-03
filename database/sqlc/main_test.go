package database

import (
	"database/sql"
	"log"
	"testing"

	"github.com/dqrk0jeste/letscube-backend/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal(err)
	}

	connection, err := sql.Open(config.DatabaseDriver, config.DatabaseSource)
	if err != nil {
		log.Fatal(err)
	}

	testQueries = New(connection)
}
