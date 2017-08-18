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

