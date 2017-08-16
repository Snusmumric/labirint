package src

import (
	"net/http"
	"fmt"
)

func moveAction(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm() // Must be called before writing response
	if err != nil {
		fmt.Fprintf(writer, anError, err)
	} else {
		response := "test"
		fmt.Fprint(writer, response)
	}
}
