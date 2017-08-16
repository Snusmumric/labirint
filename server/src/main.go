package src

import (
	"fmt"
	"log"
	"net/http"
	//"strings"
	//"strconv"
)

const (
	pageTop    = `<!DOCTYPE HTML><html align="center"><head>
<style>.error{color:#FF0000;}</style></head><title>Statistics</title>
<body><h3>Statistics</h3>
<p>Computes basic statistics for a given list of numbers</p>`
	form       = `<form action="/" method="POST">
<label for="numbers">Numbers (comma or space-separated):</label><br />

<table border="1" align="center">
<tr><td><a>Your numbers: </a></td><td><input type="text" name="numbers" size="30"></td></tr>
<tr><td><a>Answer: </a></td><td><input type="text" name="text answer" size="30" value="answer"></td></tr>
</table>

<p><input type="submit" value="Calculate"></p>

</form>`
	pageBottom = `</body></html>`
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
	fmt.Fprint(writer, pageTop, form)
	if err != nil {
		fmt.Fprintf(writer, anError, err)
	}
}


