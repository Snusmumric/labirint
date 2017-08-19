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
	http.HandleFunc("/move", moveAction)
	//http.HandleFunc("/start", startAction)
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

func moveAction(res http.ResponseWriter, req *http.Request)  { // error
	err := req.ParseForm()
	if err != nil {
		fmt.Fprintf(res, "Error with url parsing occupyed\n")
		fmt.Fprintf(res,"Error with id input occupyed\n")
		return
	}

	PlayerStr := req.URL.Query().Get("PlayerId")
	PlayerId, err := strconv.Atoi(PlayerStr)
	_ = PlayerId
	if err != nil {
		fmt.Fprintf(res, "Error with PlayerId input occupyed\n")
		fmt.Fprintf(res,"Error with id input occupyed\n")
		return
	}

	GameStr := req.URL.Query().Get("GameId")
	GameId, err := strconv.Atoi(GameStr)
	if err != nil {
		fmt.Fprintf(res, "Error with id input occupyed\n")
		fmt.Fprintf(res,"Error with id input occupyed\n")
		return
	}

	MoveStr := req.URL.Query().Get("Move")

	gameToUse := globGameBase[GameId]

	switch MoveStr {
	case "up":
		gameToUse.Gg.Pos.Posy++
	case "dn":
		gameToUse.Gg.Pos.Posy--
	case "rt":
		gameToUse.Gg.Pos.Posx++
	case "lf":
		gameToUse.Gg.Pos.Posy--
	}

	fmt.Fprintf(res, "You have entered another dungeon! %s\n", MoveStr)
}
