package main

import (
	"db_client"
	"fmt"
	"game"
	"io/ioutil"
	"log"
	"net/http"
	"reporter"
	"strconv"
	"user"
	"worker_pool"

	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type userbase map[int]user.User       // int - id of players
type globalScorebase map[int]int      // [user_id] score (global score!? == sum of scores from all games?)           !!!!!!!!
type globalGameBase map[int]game.Game //[game_id]game

const (
	html_dir         = "server/html/"
	anError          = `<p class="error"i>%s</p>`
	commonMapSize    = 5
	poolCap          = 10
	taskSetTimeoutMs = 200
)

var WP *worker_pool.WorkerPool
var LabDB *db_client.DBClient // for userList

func init() {
	WP = worker_pool.NewPool(poolCap)
	WP.Run()
	LabDB = db_client.New()
}

func main() {

	http.HandleFunc("/", homePage)
	http.HandleFunc("/move", moveAction)
	http.HandleFunc("/start", startAction)
	//http.HandleFunc("/save", saveAction)
	//http.HandleFunc("/end", endAction)
	if err := http.ListenAndServe(":9001", nil); err != nil {
		log.Fatal("failed to start server", err)
	}

	//game3, err := game.GetTheGame(3,commonMapSize,LabDB) //+
	//game3.Map_master.Field[0][0].Kind = 0 //+
	//game3.Map_master.Field[0][1].Kind = 0 //+
	//game3.Status=0 // +
	//game3.SavedName = "game4" //+
	//err = game.UpdateTheGame(game3,LabDB)//+
	//fmt.Println(err)
	//_,_ = game.GetTheGame(1,commonMapSize,LabDB)
	//_, _ = user.GetUserById(LabDB,1)

	//correct, _ := LabDB.UserPassCorrect("Dima", "123")
	//fmt.Println(correct)

}

func moveAction(writer http.ResponseWriter, request *http.Request) {
	/*
		input:
			id - player
			gameid
			mvdir - up,dn,rt,lf
		out:
			text of move event
	*/

	err := request.ParseForm()
	if err != nil {
		reporter.SendResp(writer, 400, fmt.Errorf("Invalid move request"), reporter.EmptyBody)
		return
	}

	// get userId
	str := request.URL.Query().Get("id")
	id, err := strconv.Atoi(str)
	if err != nil {
		reporter.SendResp(writer, 400, err, reporter.EmptyBody)
		return
	}

	//get gameId
	str = request.URL.Query().Get("gameid")
	gameid, err := strconv.Atoi(str)
	if err != nil {
		reporter.SendResp(writer, 400, err, reporter.EmptyBody)
		return
	}

	// user exists
	exists, err := LabDB.UserIdExists(id)
	if !exists {
		reporter.SendResp(writer, 400, fmt.Errorf("User doesn't exist: %s", err), reporter.EmptyBody)
		return
	}

	// get the user
	usr, err := user.GetUserById(LabDB, id)
	if err != nil {
		reporter.SendResp(writer, 400, fmt.Errorf("Get User Error: %s", err), reporter.EmptyBody)
		return
	}

	// game exists in user list
	gameExistsInUserList, err := usr.UserHaveTheGameWithId(gameid)
	if !gameExistsInUserList {
		reporter.SendResp(writer, 400, fmt.Errorf("This user doesn't own this game: %d", gameid), reporter.EmptyBody)
		return
	}
	if err != nil {
		reporter.SendResp(writer, 400, fmt.Errorf("Get GameExistsInUserList Error: %s", err), reporter.EmptyBody)
		return
	}

	// game exist in data base
	gameExistsInDB, err := LabDB.GameExists(gameid)
	if !gameExistsInDB {
		reporter.SendResp(writer, 400, fmt.Errorf("There is no game with such id in DataBase: %d", gameid), reporter.EmptyBody)
		return
	}
	if err != nil {
		reporter.SendResp(writer, 400, fmt.Errorf("Get GameExistsInDataBase Error: %s", err), reporter.EmptyBody)
		return
	}

	//get the game
	gameToExecute, err := game.GetTheGame(3, commonMapSize, LabDB)
	if err != nil {
		reporter.SendResp(writer, 400, fmt.Errorf("Game couldn't be loaded from DataBase: %s", err), reporter.EmptyBody)
		return
	}

	//get move direction
	MoveStr := request.URL.Query().Get("mvdir")

	//move, move, move!
	switch MoveStr {
	case "up":
		gameToExecute.Gg.Position.Posy++
	case "dn":
		gameToExecute.Gg.Position.Posy--
	case "rt":
		gameToExecute.Gg.Position.Posx++
	case "lf":
		gameToExecute.Gg.Position.Posx--
	}

	//process the event after move
	// current postion
	x := gameToExecute.Gg.Position.Posx
	y := gameToExecute.Gg.Position.Posy
	//get the event and null the room
	eventId := gameToExecute.Map_master.Field[x][y].Kind
	gameToExecute.Map_master.Field[x][y].Hidden = 0
	gameToExecute.Map_master.Field[x][y].Kind = 0

	switch {
	case eventId > 0:
		gameToExecute.Gg.Healthpoints++
		fmt.Fprintf(writer, "Smile, you get HP++.")
	case eventId < 0:
		gameToExecute.Gg.Healthpoints--
		fmt.Fprintf(writer, "It's a pity, you get injured.")
	default:
		fmt.Println("Nothing interesting here.\n Just a mermaid on the tree  seductively smiling at you...")
	}

	//update the game
	err = game.UpdateTheGame(gameToExecute, LabDB) //+
	fmt.Println(err)
}

func homePage(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm() // Must be called before writing response
	main_page_template, err := ioutil.ReadFile(html_dir + "home_page.html")
	fmt.Fprint(writer, string(main_page_template))
	if err != nil {
		fmt.Fprintf(writer, anError, err)
	}
}

// start?id=1   userId
func startAction(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()

	if err != nil {
		reporter.SendResp(writer, 400, fmt.Errorf("Invalid json request"), reporter.EmptyBody)
		return
	}

	str := request.URL.Query().Get("id")
	id, err := strconv.Atoi(str)

	if err != nil {
		reporter.SendResp(writer, 400, err, reporter.EmptyBody)
		return
	}

	usr, err := user.GetUserById(LabDB, id)

	if err != nil {
		reporter.SendResp(writer, 400, fmt.Errorf("Get User Error: %s", err), reporter.EmptyBody)
		return
	}

	notSavedGame, err := usr.CheckActiveGame(LabDB)
	if err != nil {
		reporter.SendResp(writer, 400, fmt.Errorf("CheckActiveGames error: %s", err), reporter.EmptyBody)
		return
	}

	if notSavedGame {
		reporter.SendResp(writer, 400, fmt.Errorf("not_saved_game"), reporter.EmptyBody)
		return
	}

	gs := user.GameStarter{
		DataBase: LabDB,
		User:     usr,
		MapSize:  commonMapSize,
		Writer:   writer,
	}
	err = WP.AddTaskAsync(&gs, taskSetTimeoutMs)

	if err != nil {
		reporter.SendResp(writer, 400, fmt.Errorf("timeout"), reporter.EmptyBody)
	}
}

func login(writer http.ResponseWriter, request *http.Request) {
	/*
		input: name, password
		output: id!=0 / 0 if error (not exists, wrong combination)
	*/
	err := request.ParseForm()
	if err != nil {
		reporter.SendResp(writer, 400, fmt.Errorf("Invalid login request"), reporter.EmptyBody)
		return
	}

	// get user-name
	name := request.URL.Query().Get("name")
	if name == "" {
		reporter.SendResp(writer, 400, fmt.Errorf("Name is empty."), reporter.EmptyBody)
		return
	}

	// get password
	pass := request.URL.Query().Get("pass")
	if pass == "" {
		reporter.SendResp(writer, 400, fmt.Errorf("Password is empty."), reporter.EmptyBody)
		return
	}

	// check user exists
	NameExists, err := LabDB.UserNameExists(name)
	if !NameExists {
		reporter.SendResp(writer, 400, fmt.Errorf("Such user doesn't exist: %s", err), reporter.EmptyBody)
		return
	}

	CorectPass, err := LabDB.UserPassCorrect(name, pass)
	if !CorectPass {
		reporter.SendResp(writer, 400, fmt.Errorf("Password is wrong: %s", err), reporter.EmptyBody)
		return
	}

	fmt.Println(CorectPass)

	fmt.Fprintf(writer, "%d")
}

func register(res http.ResponseWriter, req *http.Request) {
	/*
		input: name, password
		output: id!=0 / 0 if error (already exists, smth weird)
	*/

}
