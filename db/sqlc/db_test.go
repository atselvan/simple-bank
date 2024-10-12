package db

import (
	"database/sql"
	"fmt"
	"github.com/atselvan/simple-bank/util"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var (
	testQueries *Queries
	testStore   *Store
)

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	dbSource := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		config.DBUserName,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	conn, err := sql.Open(config.DBDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to the db: ", err)
	}
	testQueries = New(conn)

	testStore = NewStore(conn)

	os.Exit(m.Run())
}
