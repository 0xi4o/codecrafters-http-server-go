package main

import (
	"fmt"
	"net/http"
	"path"
)

type ResponseHeaders struct {
	ContentType   string
	ContentLength int
}

type Response struct {
	Version string
	Status  string
	Headers ResponseHeaders
	Body    string
}

func SerializeResponse(request Request) Response {
	if match, _ := path.Match("/echo/*", request.Path); match {
		body := path.Base(request.Path)
		responseHeaders := ResponseHeaders{
			ContentType:   "text/plain",
			ContentLength: len(body),
		}
		return Response{
			Version: request.Version,
			Status:  fmt.Sprintf("%d %s", http.StatusOK, "OK"),
			Headers: responseHeaders,
			Body:    body,
		}
	}

	if path.Base(request.Path) == "/" {
		return Response{
			Version: request.Version,
			Status:  fmt.Sprintf("%d %s", http.StatusOK, "OK"),
		}
	}

	return Response{
		Version: request.Version,
		Status:  fmt.Sprintf("%d %s", http.StatusNotFound, "Not Found"),
	}
}
