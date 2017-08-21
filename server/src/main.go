package main

import (
	"encoding/json"
	"fmt"
	"game"
	"gmap"
	"log"
	"net/http"
	"user"
	"worker_pool"
	//"strings"
	"io/ioutil"
	"strconv"
)

type userbase map[int]user.User       // int - id of players
type globalScorebase map[int]int      // [user_id] score (global score!? == sum of scores from all games?)           !!!!!!!!
type globalGameBase map[int]game.Game //[game_id]game

const (
	html_dir      = "server/html/"
	anError       = `<p class="error"i>%s</p>`
	commonMapSize = 5
)

type apiResp struct {
	Error error    `json: error`
	Code  int      `json: code`
	Body  struct{} `json: body`
}

func SendResp(writer http.ResponseWriter, code int, err error, body struct{}) {
	fmt.Fprintf(writer, json.Marshall(
		apiResp{
			Error: err,
			Code:  code,
			Body:  body,
		},
	),
	)
	return

}

const InvalidRequest appError = appError{
	Error: fmt.Errorf("invalid json request"),
	Code:  400,
}

const (
	poolCap = 10
)

func init() {
	WP := worker_pool.NewPool(poolCap)
	WP.Run()
}

func main() {
	http.HandleFunc("/", homePage)
	//http.HandleFunc("/move", moveAction)
	http.HandleFunc("/start", startAction)
	//http.HandleFunc("/save", saveAction)
	//http.HandleFunc("/end", endAction)
	if err := http.ListenAndServe(":9001", nil); err != nil {
		log.Fatal("failed to start server", err)
	}
}

func GetUser(user_id int) (user.User, error) {
	user, ok := userDataBase[user_id]
	if !ok {
		return nil, fmt.Errorf("Warn: no such user in base!")
	}

	return user, nil
}

func homePage(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm() // Must be called before writing response
	main_page_template, err := ioutil.ReadFile(html_dir + "home_page.html")
	fmt.Fprint(writer, string(main_page_template))
	if err != nil {
		fmt.Fprintf(writer, anError, err)
	}
}

func startAction(writer http.ResponseWriter, req *http.Request) {
	err := request.ParseForm()

	if err != nil {
		fmt.Fprintf(writer, json.Marshal(InvalidRequest))
		return
	}

	str := req.URL.Query().Get("id")
	id, err := strconv.Atoi(str)

	if err != nil {
		SendResp(writer, 400, err)
		return
	}

	user, err := GetUser(id)

	if err != nil {
		SendResp(writer, 400, fmt.Errorf("Get User Error: %s", err), struct{})
		return
	}

	notSavedGame, err := pUser.CheckActiveGame()
	if err != nil {
		SendResp(writer, 400, fmt.Errorf("CheckActiveGames error: %s", err), struct{})
		return
	}

	if notSavedGame {
		SendResp(writer, 400, fmt.Errorf("not_saved_game"), struct{})
		return
	}

	gs := user.GameStarter{}
	err = WR.AddAsynkTask(&gs)

	if err != nil {
		SendResp(writer, 400, fmt.Errrorf("timeout"), struct{})
	}
}
