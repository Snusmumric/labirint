package user

import (
	"db_client"
	"fmt"
	"game"
	"net/http"
	"reporter"
	"strings"
	"strconv"
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
	game, err := game.MakeAGame(gs.MapSize, "game1", gs.DataBase)
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
	rows, err := dbc.DB.Query("SELECT status FROM games g INNER JOIN users u ON g.id = ANY(u.games) WHERE u.id=$1", u.Id)
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
	row, err := dbc.DB.Query("SELECT * FROM users WHERE id=$1", id)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	var user User
	row.Next()

	var uint8Games []uint8
	err = row.Scan(&user.Id, &user.Name, &uint8Games)
	str := fmt.Sprintf("%s", uint8Games)
	str = str[1:len(str)-1]
	NumList := strings.Split(str, ",")
	for _, s := range NumList {
		j,_ := strconv.Atoi(s)
		user.Games = append(user.Games,j)
	}

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *User)UserHaveTheGameWithId(gameId int) (bool, error) {
	exists := false

	for _, g := range us.Games {
		if g == gameId {
			exists = true
			return exists, nil
		}
	}

	return exists, nil
}
