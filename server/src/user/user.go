package user

import (
	"db_client"
	"fmt"
	"game"
	"net/http"
	"reporter"
	"strconv"
	"strings"
)

type User struct {
	Id     int    `json:"user id number,omitempty"`
	Name   string `json:"user name,omitempty"`
	Pass   string
	Games  []int       `json:"game ids',omitempty"`     // game_id`s
	Scores map[int]int `json:"games' scores,omitempty"` // [game_id`s]score
}

type GameStarter struct {
	DataBase   *db_client.DBClient
	User       *User
	MapSize    int
	Writer     *http.ResponseWriter
	ConnHolder chan int
	gameId     int
	EventNum   int
}

type GameStarterResp struct {
	GameId int `json:"game_id"`
}

func (gs *GameStarter) Handle() error {
	fmt.Println("Handle started")
	fmt.Println(fmt.Sprintf("gs: %#v", gs))
	game, err := game.MakeAGame(gs.MapSize, "", gs.EventNum, gs.DataBase)
	if err != nil {
		return fmt.Errorf("GameStarter: failed to make a game: %s", err)
	}
	gs.User.Games = append(gs.User.Games, game.Id)
	_, err = gs.DataBase.DB.Query("UPDATE users SET games=array_append(games,$1) WHERE id=$2", game.Id, gs.User.Id)
	if err != nil {
		return fmt.Errorf("GameStarter: failed to update userDB: %s", err)
	}
	gs.gameId = game.Id
	fmt.Println("Handle finished")
	return nil
}

func (gs *GameStarter) Finish(err error) {
	fmt.Println("Finish starter")
	if err != nil {
		fmt.Fprintf(*gs.Writer, "GameStarter failed: %s", err)
		gs.ConnHolder <- 0
	}
	body := GameStarterResp{
		GameId: gs.gameId,
	}
	reporter.SendResp(*gs.Writer, 200, nil, body)

	gs.ConnHolder <- 0
	fmt.Println("Finish finished")
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
	if !row.Next() {
		return nil, fmt.Errorf("not_exists")
	}

	var uint8Games []uint8
	err = row.Scan(&user.Id, &user.Name, &user.Pass, &uint8Games)
	if err != nil {
		return nil, err
	}
	str := fmt.Sprintf("%s", uint8Games)
	str = str[1 : len(str)-1]
	NumList := strings.Split(str, ",")
	for _, s := range NumList {
		j, _ := strconv.Atoi(s)
		user.Games = append(user.Games, j)
	}

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetIdByUserName(name string, dbc *db_client.DBClient) (int, error) {

	strToExec := fmt.Sprintf("SELECT id FROM users WHERE name='%s'", name)
	row, err := dbc.DB.Query(strToExec)
	fmt.Println(err)
	if err != nil {
		return 0, err
	}
	defer row.Close()

	row.Next()
	var id int
	err = row.Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (us *User) UserHaveTheGameWithId(gameId int) (bool, error) {
	exists := false

	for _, g := range us.Games {
		if g == gameId {
			exists = true
			return exists, nil
		}
	}

	return exists, nil
}
