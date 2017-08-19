package main

import (
	"fmt"
	"log"
	"net/http"
	//"strings"
	//"strconv"
	"io/ioutil"
)

const (
	html_dir = "server/html/"
	anError    = `<p class="error">%s</p>`
)



func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/start", starter)
	http.HandleFunc("/move", moveAction)
	http.HandleFunc("/status", showStatus)
	http.HandleFunc("/login", login)
	//http.HandleFunc("/save", saveAction)
	//http.HandleFunc("/end", endAction)
	if err := http.ListenAndServe(":9001", nil); err != nil {
		log.Fatal("failed to start server", err)
	}
}

func homePage(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm() // Must be called before writing response
	main_page_template, err := ioutil.ReadFile(html_dir+"home_page.html")
	fmt.Fprint(writer, string(main_page_template))
	if err != nil {
		fmt.Fprintf(writer, anError, err)
	}
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

	namestr := req.URL.Query().Get("name")

	// ckeck the existance of the user!
	if olduser, ok := userDataBase[namestr]; ok {
		fmt.Fprintf(res, "Hello again, %s!\n", namestr)
		fmt.Fprintf(res, "We have created a new game for you now :)\n")
		globalGameNum++
		newgame := game.MakeAGame(3, len(globGameBase)+1)
		olduser.Games = append(olduser.Games, newgame.Id)
		olduser.Scores[newgame.Id] = 0 // initialize score in user scoretable
		userDataBase[olduser.Name] = olduser
		globGameBase[newgame.Id] = *newgame
	} else {
		fmt.Fprintf(res, "Hello, newbie!\n")
		fmt.Fprintf(res, "Let`s have a dungeon ride!\n")
		newuser := user.User{Id: id, Name: namestr, Scores: map[int]int{}}
		globalGameNum++
		newgame := game.MakeAGame(3, len(globGameBase)+1)
		newuser.Games = append(newuser.Games, newgame.Id)
		newuser.Scores[newgame.Id] = 0 // initialize score in user scoretable

		// updating global databases:
		userDataBase[newuser.Name] = newuser
		userScoreDataBase[newuser.Id] = 0
		globGameBase[newgame.Id] = *newgame

		//newmap :=
		userNum++
		respdata, _ := json.Marshal(newuser)
		fmt.Fprintf(res, "New user constructed: %s\n", respdata)
		fmt.Printf("%s\n", respdata)
	}

}

func login(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		fmt.Fprintf(res, "Error with url parsing while login occupyed\n")
		return
	}

	str := req.URL.Query().Get("name")
	if str == "" {
		fmt.Fprintf(res, "Error: your name is empty!\n")
		return
	}

	// look for a name in the base
	if val, ok := userDataBase[str]; ok {
		fmt.Fprintf(res, "Hello again, %s!\n", str)
		respdata, _ := json.Marshal(val)
		fmt.Fprintf(res, "Here are you saved games: %s\n", respdata)
		fmt.Fprintf(res, "\nWould you like to continue any of them of try a new one?\n")
	} else {
		newId := len(userDataBase) + 1
		userNum++
		newuser := user.User{Id: newId, Name: str, Scores: map[int]int{}}
		globalGameNum++
		newgame := game.MakeAGame(3, len(globGameBase)+1)
		newuser.Games = append(newuser.Games, newgame.Id)
		newuser.Scores[newgame.Id] = 0
		userScoreDataBase[newuser.Id] = 0
		globGameBase[newgame.Id] = *newgame
		userDataBase[newuser.Name] = newuser
		fmt.Fprintf(res, "We have made a new account for you, %s.\n", str)

		respdata, _ := json.Marshal(newuser)
		fmt.Fprintf(res, "New user constructed: %s\n", respdata)
		fmt.Printf("%s\n", respdata)
	}



}

func moveAction(res http.ResponseWriter, req *http.Request) { // error
	err := req.ParseForm()
	if err != nil {
		fmt.Fprintf(res, "Error with url parsing occupyed\n")
		//return fmt.Errorf("Error with url parsing occupyed\n")
		fmt.Fprintf(res, "Error with id input occupyed\n")
		return
	}

	PlayerStr := req.URL.Query().Get("PlayerId")
	PlayerId, err := strconv.Atoi(PlayerStr)
	_ = PlayerId
	if err != nil {
		fmt.Fprintf(res, "Error with PlayerId input occupyed\n")
		//return fmt.Errorf("Error with PlayerId input occupyed\n")
		fmt.Fprintf(res, "Error with id input occupyed\n")
		return
	}

	GameStr := req.URL.Query().Get("GameId")
	GameId, err := strconv.Atoi(GameStr)
	if err != nil {
		fmt.Fprintf(res, "Error with id input occupyed\n")
		//return fmt.Errorf("Error with id input occupyed\n")
		fmt.Fprintf(res, "Error with id input occupyed\n")
		return
	}

	MoveStr := req.URL.Query().Get("Move")

	gameToMod := globGameBase[GameId]

	switch MoveStr {
	case "up":
		gameToMod.Gg.Pos.Posy++
	case "dn":
		gameToMod.Gg.Pos.Posy--
	case "rt":
		gameToMod.Gg.Pos.Posx++
	case "lf":
		gameToMod.Gg.Pos.Posy--
	}

	globGameBase[GameId] = gameToMod

	fmt.Fprintf(res, "You have entered another dungeon! %s\n", MoveStr)
	//return nil
}

func showStatus(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		fmt.Fprintf(res, "Error with url parsing occupyed\n")
		//return fmt.Errorf("Error with url parsing occupyed\n")
		fmt.Fprintf(res, "Error with id input occupyed\n")
		return
	}

	/*
	PlayerStr := req.URL.Query().Get("PlayerId")
	PlayerId, err := strconv.Atoi(PlayerStr)
	if err != nil {
		fmt.Fprintf(res, "Error with PlayerId input occupyed\n")
		//return fmt.Errorf("Error with PlayerId input occupyed\n")
		fmt.Fprintf(res, "Error with id input occupyed\n")
		return
	}
	*/
	PlayerName := req.URL.Query().Get("PlayerName")


	GameStr := req.URL.Query().Get("GameId")
	GameId, err := strconv.Atoi(GameStr)
	if err != nil {
		fmt.Fprintf(res, "Error with id input occupyed\n")
		//return fmt.Errorf("Error with id input occupyed\n")
		fmt.Fprintf(res, "Error with id input occupyed\n")
		return
	}

	gamer := userDataBase[PlayerName]
	gameToUse := globGameBase[GameId]

	Posx := gameToUse.Gg.Pos.Posx
	Posy := gameToUse.Gg.Pos.Posy
	fmt.Fprintf(res, "%s, you are here [%d], [%d]. This room has kind of: %d.\n",
		gamer.Name, Posx, Posy, gameToUse.Field.Field[Posx][Posy].Kind)
}
