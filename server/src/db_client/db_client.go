package db_client

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	password = "pswd"
	user     = "labirint1"
	dbname   = "labirintDB"
)

type DBClient struct {
	DB *sql.DB
}

func (dbc *DBClient) init() {
	PsqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	database, err := sql.Open("postgres", PsqlInfo)
	if err != nil {
		panic(err)
	}
	_, err = database.Exec("CREATE TABLE IF NOT EXISTS games (" +
		"id SERIAL PRIMARY KEY, " +
		"map text[][] " +
		"status int" +
		"saved_name text")
	if err != nil {
		panic(err)
	}
	_, err = database.Exec("CREATE TABLE IF NOT EXISTS users (" +
		"id SERIAL PRIMARY KEY, " +
		"name TEXT, " +
		"games int[])")
	if err != nil {
		panic(err)
	}

	dbc.DB = database
	return
}

func New() *DBClient {
	dbc := &DBClient{}
	dbc.init()
	return dbc

}
