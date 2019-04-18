package main

import (
	"fmt"
	_ "./userblock"
	"./DBmanipul"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "postgres"
)

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db := DBmanipul.InitDB(psqlInfo)
	defer db.Close()

	db2name := "lobirint"
	/*
	DBmanipul.CreateDB(db, db2name)
	db2 := DBmanipul.OpenDB(db2name)

	testunit1 := userblock.User{1, "Dima"}
	userblock.CreateUserTable(db2)
	userblock.InsertUserIntoUserTable(db2, testunit1)
	testunit2 := userblock.User{1, "Kira"}
	userblock.UpdateUserNameInUserTable(db2, testunit2)
	*/

	DBmanipul.DropDB(db,db2name)
  
  }
