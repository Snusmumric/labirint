package reporter

import (
	"enconding/json"
	"fmt"
	"net/http"
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

const InvalidRequest = fmt.Errorf("invalid json request")
