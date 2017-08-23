package user

import (
	"db_client"
	"fmt"
	"game"
	"net/http"
	"reporter"
)

type User struct {
	Id     int         `json:"user id number,omitempty"`
	Name   string      `json:"user name,omitempty"`
	Games  []int       `json:"game ids',omitempty"`     // game_id`s
	Scores map[int]int `json:"games' scores,omitempty"` // [game_id`s]score
}

type GameStarter struct {
	DataBase *db_client.DBClient
	User     *User
	MapSize  int
	Writer   http.ResponseWriter
	gameId   int
}

func (gs *GameStarter) Handle() error {
	game, err := game.MakeAGame(gs.MapSize, gs.DataBase)
	if err != nil {
		return fmt.Errorf("GameStarter: failed to make a game: %s", err)
	}
	gs.User.Games = append(gs.User.Games, game.Id)
	_, err = gs.DataBase.DB.Query("UPDATE users SET games array_append(games,?) WHERE id=?", game.Id, gs.User.Id)
	if err != nil {
		return fmt.Errorf("GameStarter: failed to update userDB: %s", err)
	}
	gs.gameId = game.Id
	return nil
}

func (gs *GameStarter) Finish(err error) {
	if err != nil {
		fmt.Fprintf(gs.Writer, "GameStarter failed: %s", err)
	}
	reporter.SendResp(gs.Writer, 200, nil, struct {
		gameId int `json: game_id`
	}{gs.gameId})

}

func (u *User) CheckActiveGame(dbc *db_client.DBClient) (bool, error) {
	var count int
	rows, err := dbc.DB.Query("SELECT status FROM games g INNER JOIN users u ON g.id = ANY(u.games) WHERE u.id=?", u.Id)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		var status int
		err = rows.Scan(&status)
		if err != nil {
			return false, err
		}
		if status == 1 {
			count += 1
		}
	}
	if count > 1 {
		return false, fmt.Errorf("Impossible situation: active games number: %d", count)
	}
	if count == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func GetUserById(dbc *db_client.DBClient, id int) (*User, error) {
	row, err := dbc.DB.Query("SELECT * FROM users WHERE id=?", id)
	defer row.Close()
	if err != nil {
		return nil, err
	}
	var user User
	err = row.Scan(&user.Id, &user.Name, &user.Games)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
