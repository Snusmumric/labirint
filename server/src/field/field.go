package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type userbase map[int]user       // int - id of players
type globalScorebase map[int]int // [user_id] score (global score!? == sum of scores from all games?)           !!!!!!!!
type globalGameBase map[int]game //[game_id]game

type user struct {
	Id     int         `json:"user id number,omitempty"`
	Name   string      `json:"user name,omitempty"`
	Games  []int       `json:"game ids',omitempty"`     // game_id`s
	Scores map[int]int `json:"games' scores,omitempty"` // [game_id`s]score
}

type cell struct {
	Kind int
	/*
		>0 heal
		=0 noting or seen
		<0 damage
	*/

	Hidden int // =1(hidden) || =0(seen)
}

type playchar struct {
	Healthpoints int
	//position *cell //or i,j int
	// address of cell may not survive save-load process
	Posx int
	Posy int
}

type game struct {
	Id         int `json:"game-id,omitempty"`
	Map_master gmap
	Map_gg     gmap
	Gg         playchar
	Status     int // -1(over), 0(saved), 1(online)
}

func make_a_game() *game {
	globalGameNum++
	newmap := make_a_map(mapSize)
	newCharacter := playchar{100, 0, 0}
	newgame := game{globalGameNum, *newmap, *newmap, newCharacter, 1}
	return &newgame
}

type gmap struct {
	Field [][]cell
}

func make_a_map(size int) *gmap {
	cellarray := [][]cell{}
	for j := 0; j < size; j++ {
		cellrow := []cell{}
		for i := 0; i < size; i++ {
			celltoadd := cell{1, 1}
			cellrow = append(cellrow, celltoadd)
		}
		cellarray = append(cellarray, cellrow)
	}
	return &gmap{cellarray}
}

//=============global=========================
//-----------data bases-----------------------
var userDataBase userbase
var userScoreDataBase globalScorebase
var globGameBase globalGameBase

//--------------------------------------------
var userNum int
var globalGameNum int
var mu sync.Mutex
var countOfRequests int
var mapSize = 3

//============================================

func main() {
	// initializing global data bases:
	userDataBase = userbase{}
	userScoreDataBase = globalScorebase{}
	globGameBase = globalGameBase{}

	http.HandleFunc("/", home) // перехватывает вообще все запросы
	http.HandleFunc("/count", counter)
	http.HandleFunc("/start", starter)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func home(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	countOfRequests++
	mu.Unlock()
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

func counter(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n", countOfRequests)
	mu.Unlock()
}

func starter(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		fmt.Fprintf(res, "Error with url parsing occupyed\n")
	}

	str := req.URL.Query().Get("id")
	id, err := strconv.Atoi(str)
	if err != nil {
		fmt.Fprintf(res, "Error with id input occupyed\n")
	}

	str = req.URL.Query().Get("name")

	newuser := user{Id: id, Name: str, Scores: map[int]int{}}
	newgame := make_a_game()
	newuser.Games = append(newuser.Games, newgame.Id)
	newuser.Scores[newgame.Id] = 0 // initialize score in user scoretable

	// updating global databases:
	userDataBase[newuser.Id] = newuser
	userScoreDataBase[newuser.Id] = 0
	globGameBase[newuser.Id] = *newgame

	//newmap :=
	userNum++
	respdata, _ := json.Marshal(newuser)
	fmt.Fprintf(res, "New user constructed: %s\n", respdata)
	fmt.Printf("%s\n", respdata)

}
