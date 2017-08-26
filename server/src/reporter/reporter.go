package reporter

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type apiResp struct {
	Error string      `json: error`
	Code  int         `json: code`
	Body  interface{} `json: body`
}

func SendResp(writer http.ResponseWriter, code int, err error, body interface{}) {
	var e string
	if err == nil {
		e = ""
	} else {
		e = err.Error()
	}
	resp, _ := json.Marshal(
		apiResp{
			Error: e,
			Code:  code,
			Body:  body,
		},
	)
	fmt.Fprintf(writer, string(resp))
	return

}

const EmptyBody = ""
