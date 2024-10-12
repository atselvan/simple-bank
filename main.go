package main

import (
	"database/sql"
	"fmt"
	"github.com/atselvan/simple-bank/api"
	db "github.com/atselvan/simple-bank/db/sqlc"
	"github.com/atselvan/simple-bank/util"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
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

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
