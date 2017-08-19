package user

import (
	"fmt"

)

type User struct {
	Id int `json:"user id number,omitempty"`
	Name string `json:"user name,omitempty"`
	Games []int `json:"game ids',omitempty"` // game_id`s
	Scores map[int]int `json:"games' scores,omitempty"` // [game_id`s]score
}

func (u *User) CheckActiveGames (bool error) {
	var count int
	for _, game := u.Games {
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
