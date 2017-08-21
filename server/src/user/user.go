package user

import (
	"database/sql"
	"fmt"
	"game"
	"gmap"
	"math"
	"net/http"
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
	writer   http.ResponceWriter
}

func (gs *GameStarter) Handle() error {
	game, err := game.MakeAGame(gs.MapSizei, &gs.DataBase)
	if err != nil {
		return fmt.Errorf("GameStarter: failed to make a game: %s", err)
	}
	gs.User.Games = append(gs.User.Games, game.Id)
	_, err = gs.DataBase.Query("UPDATE user SET games array_append(games,?) WHERE id=?", game.Id, User.Id)
	if err != nil {
		return fmt.Errorf("GameStarter: failed to update userDB: %s", err)
	}
	return nil
}

func (gs *GameStarter) Finish(err error) {
	if err != nil {
		fmt.Fprintf(gs.writer, "GameStarter failed: %s", err)
	}
	// Здесь будем отдавать пользоавтелю game id в стандартной структуре ответа{code: 200, body: {game_id: 23 } }, когда вынесем ёё из main
}

func (u *User) CheckActiveGame() (bool, error) {
	var count int
	for _, game := range u.Games {
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

func StartGame() *game.Game {
	id := math.rand(100)
	game := game.MakeAGame(mapSize)
}
