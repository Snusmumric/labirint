package user

import (
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
	// user DB to write and the end Games DB
	User    User
	MapSize int
	writer  http.ResponceWriter
}

func (gs *GameStarter) Handle() error {
	gameId = math.rand(100) // Здесь фактически будет insert в базу которая сама отдаст тновй уникальный id
	game, err := game.MakeAGame(gs.MapSize, gameID)
	if err != nil {
		return fmt.Errorf("")
	}
	gs.User.Games = append(gs.User.Games, gameId)
	// put into user DB gameid
	// put into game db
	// put into memcache
	return
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
