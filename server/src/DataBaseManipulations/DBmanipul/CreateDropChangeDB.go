package DBmanipul

import (
	"database/sql"
	"fmt"
)


const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "postgres"
)

func InitDB(psqlInfo string) *sql.DB {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}
	return db
}

// creating db ++
func CreateDB(db *sql.DB, newDBName string) {
	var strToExec string
	strToExec = fmt.Sprintf(`
	CREATE DATABASE %s;
	`, newDBName)

	_, err := db.Exec(strToExec)
	if err != nil {
		panic(err)
	}
}

// creating db ++
func DropDB(db *sql.DB, newDBName string) {
	var strToExec string
	strToExec = fmt.Sprintf(`
	DROP DATABASE %s;
	`, newDBName)

	_, err := db.Exec(strToExec)
	if err != nil {
		panic(err)
	}
}

// trying to move to another db by postgres commands, hopeless
// down here is another func to do this by db.Open
// '\' ?! in pq
// won't work
func ChangeTheDB(db *sql.DB, moveToDB string) *sql.DB {

	var strToExec string
	strToExec = fmt.Sprintf("\\c %s", moveToDB)

	_, err := db.Exec(strToExec)
	if err != nil {
		panic(err)
	}
	return db
}

// open another db by inputed the DBname ++
func OpenDB(DBname string) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, DBname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}
	return db
}

// create table if not exists ++
func CreateTable(db *sql.DB, tableName string) {

	var strToExec string
	strToExec = fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s(
		Id INT,
		Name TEXT
	);
	`, tableName)

	_, err := db.Exec(strToExec)
	if err != nil {
		panic(err)
	}
}

