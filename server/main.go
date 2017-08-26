package main

import (
	"db_client"
	"encoding/json"
	"fmt"
	"game"
	"io/ioutil"
	"log"
	"net/http"
	"reporter"
	"strconv"
	"user"
	"worker_pool"
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
var eventTextDB = make(map[int]string)
var eventBonusDB = make(map[int]int)

func init() {
	WP = worker_pool.NewPool(poolCap)
	WP.Run()
	LabDB = db_client.New()

	eventTextDB[0] = "Nothing interesting here.\n Just a mermaid on the tree  seductively smiling at you..."
	eventBonusDB[0] = 0

	eventTextDB[1] = "Smile, a mermaid have kissed you. HP++"
	eventBonusDB[1] = 1

	eventTextDB[2] = "It's a pity. Troll have injured you with a punch. HP--"
	eventBonusDB[2] = -1

	eventTextDB[3] = "Heal portion was just under your leg!"
	eventBonusDB[3] = 1
}

func main() {

	http.HandleFunc("/", homePage)
	http.HandleFunc("/start", startAction)

	http.HandleFunc("/login", loginAction)
	http.HandleFunc("/register", registerAction)
	http.HandleFunc("/move", moveAction)
	http.HandleFunc("/newGame", newGameAction)

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
	//UserId, err := user.GetIdByUserName("Dim", LabDB)

	//game3, _ := game.MakeAGame(5, "game1", LabDB)
	//fmt.Println(game3.Id)
	//id, err := LabDB.RegistrateNewUser("Kira", "234")
	//fmt.Println(err)
	//fmt.Println(id)

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
	gameToExecute, err := game.GetTheGame(gameid, commonMapSize, LabDB)
	if err != nil {
		reporter.SendResp(writer, 400, fmt.Errorf("Game couldn't be loaded from DataBase: %s", err), reporter.EmptyBody)
		return
	}

	//get move direction
	MoveStr := request.URL.Query().Get("mvdir")

	//move, move, move!
	switch MoveStr {
	case "up":
		if gameToExecute.Gg.Position.Posy < commonMapSize-1 {
			gameToExecute.Gg.Position.Posy++
		}
	case "down":
		if gameToExecute.Gg.Position.Posy > 0 {
			gameToExecute.Gg.Position.Posy--
		}
	case "right":
		if gameToExecute.Gg.Position.Posx < commonMapSize-1 {
			gameToExecute.Gg.Position.Posx++
		}
	case "left":
		if gameToExecute.Gg.Position.Posy > 0 {
			gameToExecute.Gg.Position.Posx--
		}
	case "stand":

	}

	//process the event after move
	// current postion
	x := gameToExecute.Gg.Position.Posx
	y := gameToExecute.Gg.Position.Posy

	//get the event
	eventId := gameToExecute.Map_master.Field[x][y].Kind
	eventToAppearT := eventTextDB[eventId]
	eventToAppearB := eventBonusDB[eventId]

	//null the room
	gameToExecute.Map_master.Field[x][y].Hidden = 0
	gameToExecute.Map_master.Field[x][y].Kind = 0

	gameToExecute.Gg.Healthpoints += eventToAppearB

	//update the game
	gameToExecute.MapEventRandomizator(len(eventTextDB) - 1)
	err = game.UpdateTheGame(gameToExecute, LabDB) // update in DB//+
	if err != nil {
		fmt.Println(err)
		return
	}

	c := struct {
		User_id int    `json:"user_id"`
		Game_id int    `json:"game_id"`
		X       int    `json:"x"`
		Y       int    `json:"y"`
		Event   string `json:"event"`
		Bonus   int    `json:"bonus"`
	}{User_id: usr.Id, Game_id: gameToExecute.Id, X: gameToExecute.Gg.Position.Posx, Y: gameToExecute.Gg.Position.Posy, Event: eventToAppearT, Bonus: eventToAppearB}
	jdata, err := json.Marshal(&c)
	if err != nil {
		fmt.Println("json error in move action:", err)
	}
	//os.Stdout.Write(jdata)
	fmt.Fprintf(writer, "%s", jdata)
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

	fmt.Println(fmt.Sprintf("wrinter: %v", &writer))
	holder := make(chan int)
	gs := user.GameStarter{
		DataBase:   LabDB,
		User:       usr,
		MapSize:    commonMapSize,
		Writer:     &writer,
		ConnHolder: holder,
		EventNum:   len(eventTextDB),
	}
	err = WP.AddTaskAsync(&gs, taskSetTimeoutMs)

	_ = <-holder

	if err != nil {
		reporter.SendResp(writer, 400, fmt.Errorf("timeout"), reporter.EmptyBody)
	}
}

//+
func loginAction(writer http.ResponseWriter, request *http.Request) {
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
		reporter.SendResp(writer, 400, fmt.Errorf("name_required"), reporter.EmptyBody)
		return
	}

	// get password
	pass := request.URL.Query().Get("pass")
	if pass == "" {
		reporter.SendResp(writer, 400, fmt.Errorf("password_required"), reporter.EmptyBody)
		return
	}

	// check user exists
	NameExists, err := LabDB.UserNameExists(name)
	if !NameExists {
		reporter.SendResp(writer, 400, fmt.Errorf("not_exist"), reporter.EmptyBody)
		return
	}

	CorectPass, err := LabDB.UserPassCorrect(name, pass)
	if !CorectPass {
		reporter.SendResp(writer, 400, fmt.Errorf("password_wrong"), reporter.EmptyBody)
		return
	}

	UserId, err := user.GetIdByUserName(name, LabDB)
	if err != nil {
		reporter.SendResp(writer, 400, fmt.Errorf("Some problem in login logic!\nTell the admin!"), reporter.EmptyBody)
		return
	}

	fmt.Fprintf(writer, "%d", UserId)
}

//+
func registerAction(writer http.ResponseWriter, request *http.Request) {
	/*
		input: name, password
		output: id!=0 / 0 if error (already exists, smth weird)
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
	if NameExists {
		reporter.SendResp(writer, 400, fmt.Errorf("Such user slready exists!"), reporter.EmptyBody)
		return
	}

	// registration! +
	// return id of just registrated person
	UserId, err := LabDB.RegistrateNewUser(name, pass)
	if err != nil {
		reporter.SendResp(writer, 400, fmt.Errorf("Some problem in register logic!\nTell the admin!"), reporter.EmptyBody)
		return
	}

	fmt.Fprintf(writer, "%d", UserId)

}

func newGameAction(writer http.ResponseWriter, request *http.Request) {
	/*
		input: userid
		output: newgameid!=0
	*/
	err := request.ParseForm()
	if err != nil {
		reporter.SendResp(writer, 400, fmt.Errorf("Invalid login request"), reporter.EmptyBody)
		return
	}

}
