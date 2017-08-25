package db_client

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	password = "123456"
	user     = "postgres"
	dbname   = "labirintdb"
)

type DBClient struct {
	DB *sql.DB
}

// +
func (dbc *DBClient)UserIdExists(user_id int) (bool, error) {

	found := false

	rows, err := dbc.DB.Query("select exists(select * from users where id=$1)",user_id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(&found); err != nil {
			log.Fatal(err)
			return false, err
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
			return false, err
		}
	}

	return found, nil
}

// +
func (dbc *DBClient)UserNameExists(nameToFound string) (bool, error) {

	found := false
	strToExec := fmt.Sprintf("select exists(select * from users where name='%s')",nameToFound)
	rows, err := dbc.DB.Query(strToExec)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(&found); err != nil {
			log.Fatal(err)
			return false, err
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
			return false, err
		}
	}

	return found, nil
}

func (dbc *DBClient)UserPassCorrect(name string, pass string) (bool, error) {

	found := false
	strToExect := fmt.Sprintf()
	rows, err := dbc.DB.Query("select exists(select * from users where name='$1', pass='$2' )",name, pass)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(&found); err != nil {
			log.Fatal(err)
			return false, err
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
			return false, err
		}
	}

	return found, nil
}

func (dbc *DBClient)GameExists(game_id int) (bool, error) {

	found := false

	rows, err := dbc.DB.Query("select exists(select * from games where id=$1)",game_id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(&found); err != nil {
			log.Fatal(err)
			return false, err
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
			return false, err
		}
	}

	return found, nil
}

func (dbc *DBClient) init() {
	PsqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	database, err := sql.Open("postgres", PsqlInfo)
	if err != nil {
		panic(err)
	}
	_, err = database.Exec("CREATE TABLE IF NOT EXISTS games (" +
		"id SERIAL PRIMARY KEY, " +
		"map text[][]," +
		"playchar text," +
		"status int," +
		"saved_name text)")
	if err != nil {
		panic(err)
	}
	_, err = database.Exec("CREATE TABLE IF NOT EXISTS users (" +
		"id SERIAL PRIMARY KEY, " +
		"name TEXT, " +
		"pass TEXT, " +
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
