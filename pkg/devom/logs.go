package devom

import (
	"log"
	"net/http"
)

func logRequestError(req *http.Request, err error) {
	log.Printf("[%s] 😱 %s\nerror:%s\n", req.Method, req.URL, err.Error())
}

func logRequest(req *http.Request) {
	log.Printf("[%s] %s \n", req.Method, req.URL)
}

func logRequestResponseError(req *http.Request, resp *http.Response, body []byte, err error) {
	log.Printf(
		"[%s] 😱 %s \n%s\npayload: %s\nreponse status: %d",
		req.Method,
		req.URL,
		err.Error(),
		string(body),
		resp.StatusCode)
}
