package game

import (
	"../gmap"
	"../playchar"
	"db_client"
)

type Game struct {
	Id         int `json:"game-id,omitempty"`
	Map_master gmap.Gmap
	Gg         playchar.Playchar
	Status     int // -1(over), 0(saved), 1(online)
	SavedName  string
}

func MakeAGame(mapSize int, db *db_client.DBCient) (*Game, error) {
	//globalGameNum++
	newmap := gmap.MakeAMap(mapSize)
	var id int
	err := db.Query("INSERT INTO games (status map) VALUES (? ?) RETURNING id", "online", newmap.InsertString()).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("MakeAGame: failed to insert into games %s", err)
	}

	newCharacter := playchar{100, 0, 0}

	newgame := Game{id, *newmap, newCharacter, 1}

	return &newgame, nil
}
