package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://sbAdmin:postgres@localhost:5432/simple-bank?sslmode=disable"
)

var (
	testQueries *Queries
	testStore   *Store
)

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to the db: ", err)
	}
	testQueries = New(conn)

	testStore = NewStore(conn)

	os.Exit(m.Run())
}
