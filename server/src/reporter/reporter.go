package reporter

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type apiResp struct {
	Error error       `json: error`
	Code  int         `json: code`
	Body  interface{} `json: body`
}

func SendResp(writer http.ResponseWriter, code int, err error, body interface{}) {
	resp, _ := json.Marshal(
		apiResp{
			Error: err,
			Code:  code,
			Body:  body,
		},
	)
	fmt.Fprintf(writer, string(resp))
	return

}

const EmptyBody = ""
