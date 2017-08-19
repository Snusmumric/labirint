package game

import (
	"../gmap"
	"../playchar"
)

type Game struct {
	Id         int `json:"game-id,omitempty"`
	Map_master gmap.Gmap
	Gg         playchar.Playchar
	Status     int // -1(over), 0(saved), 1(online)
	SavedName  string
}

func MakeAGame(mapSize int, gameId int) (*Game, error) {
	//globalGameNum++
	newmap := gmap.MakeAMap(mapSize)
	newCharacter := playchar{100, 0, 0}
	newgame := Game{gameId, *newmap, newCharacter, 1}
	return &newgame, nil
}
