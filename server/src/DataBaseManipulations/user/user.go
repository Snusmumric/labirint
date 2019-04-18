package userblock

import (
	"database/sql"
	"fmt"
)

type User struct {
	Id     int `json:"user id number,omitempty"`
	Name   string `json:"user name,omitempty"`
//	Games  []int `json:"game ids',omitempty"`           // game_id`s
}

func CreateUserTable(db *sql.DB) {

	tableName := "Users"

	var strToExec string
	strToExec = fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s(
		Id INT unique,
		Name TEXT
	);
	`, tableName)

	_, err := db.Exec(strToExec)
	if err != nil {
		panic(err)
	}
}

func InsertUserIntoUserTable(db *sql.DB, us User) {

	tableName := "Users"

	var strToExec string
	strToExec = fmt.Sprintf(`
	insert into %s(Id, Name)
	Values (%d, '%s');
	`, tableName, us.Id, us.Name)

	_, err := db.Exec(strToExec)
	if err != nil {
		panic(err)
	}
}

func UpdateUserNameInUserTable(db *sql.DB, us User) {

	tableName := "Users"

	var strToExec string
	strToExec = fmt.Sprintf(`
	Update %s
	set name='%s'
	where id=%d;
	`, tableName,us.Name, us.Id)

	_, err := db.Exec(strToExec)
	if err != nil {
		panic(err)
	}
}
