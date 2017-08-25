package game

import (
	"db_client"
	"fmt"
	"gmap"
	"playchar"
)

type Game struct {
	Id         int `json:"game-id,omitempty"`
	Map_master *gmap.Gmap
	Gg         *playchar.Playchar
	Status     int // -1(over), 0(saved), 1(online)
	SavedName  string
}

func MakeAGame(mapSize int, dbc *db_client.DBClient) (*Game, error) {
	//globalGameNum++
	newmap := gmap.MakeAMap(mapSize)
	var id int
	strtoexec:=fmt.Sprintf("INSERT INTO games (status,map) VALUES (%s,%s) RETURNING id", "'online'", newmap.InsertString())
	res, err := dbc.DB.Query(strtoexec)
	defer res.Close()
	if err != nil {
		return nil, fmt.Errorf("MakeAGame: failed to insert into games %s", err)
	}
	err = res.Scan(&id)

	newCharacter := playchar.New(100, 0, 0)

	newgame := Game{id, newmap, newCharacter, 1, ""}

	return &newgame, nil
}
