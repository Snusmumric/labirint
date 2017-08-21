package db_client

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "labirint"
	password = "StrongPassword1"
	dbname   = "labirintDB"
)

type DBClient struct {
	db *sql.DB
}

func (udb *DBClient) init() {
	PsqlInfo := fmt.Printf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	database, err := sql.Open("postgres", PsqlInfo)
	if err != nil {
		panic(err)
	}
	_, err := database.Exec("CREATE TABLE IF NOT EXISTS games (" +
		"id SERIAL PRIMARY KEY, " +
		"map text[][] " +
		"status text" +
		"saved_name text")
	if err != nil {
		panic(err)
	}
	_, err := database.Exec("CREATE TABLE IF NOT EXISTS users (" +
		"id SERIAL PRIMARY KEY, " +
		"name TEXT, " +
		"games integer[])")
	if err != nil {
		panic(err)
	}

	gdb.db = database
	return
}

func New() *GameDBClient {
	gdb = GameDBClient
	gdb.Init()
	return gdb

}
