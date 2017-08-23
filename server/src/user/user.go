package user

import (
	"database/sql"
	"db_client"
	"fmt"
	"game"
	"gmap"
	"math"
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
	DataBase db_client.DBClient
	User     User
	MapSize  int
	Writer   http.ResponceWriter
	gameId   int
}

func (gs *GameStarter) Handle() error {
	game, err := game.MakeAGame(gs.MapSizei, &gs.DataBase)
	if err != nil {
		return fmt.Errorf("GameStarter: failed to make a game: %s", err)
	}
	gs.User.Games = append(gs.User.Games, game.Id)
	_, err = gs.DataBase.db.Query("UPDATE users SET games array_append(games,?) WHERE id=?", game.Id, User.Id)
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
	reporter.SendResp(gs.Writer, 200, "", struct {
		game_id int
	}{game.Id})

}

func (u *User) CheckActiveGame(dbc *DBClient) error {
	var count int
	var games []int
	rows, err := dbc.db.Query("SELECT status FROM games g INNER JOIN users u ON g.id = ANY(u.games) WHERE u.id=?", u.Id)
	if err != nil {
		return err
	}
	defer rows.Close()
	err = rows.Scan(&games)
	if err != nil {
		return err
	}
	for _, game := range games {
		if game.Status == 1 {
			count += 1
		}
	}
	if count > 1 {
		return 0, fmt.Errorf("Impossible situation: active games number: %d", count)
	}
	if count == 1 {
		return true, nil
	} else {
		return false, nil
	}
}
