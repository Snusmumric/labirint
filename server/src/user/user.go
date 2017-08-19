package user

import (

)

type User struct {
	Id int `json:"user id number,omitempty"`
	Name string `json:"user name,omitempty"`
	Games []int `json:"game ids',omitempty"` // game_id`s
	Scores map[int]int `json:"games' scores,omitempty"` // [game_id`s]score
}
