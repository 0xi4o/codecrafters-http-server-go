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

func SerializeResponse(request Request) string {
	if match, _ := path.Match("/user-agent", request.Path); match {
		body := request.Headers.UserAgent
		responseHeaders := ResponseHeaders{
			ContentType:   "text/plain",
			ContentLength: len(body),
		}
		response := Response{
			Version: request.Version,
			Status:  fmt.Sprintf("%d %s", http.StatusOK, "OK"),
			Headers: responseHeaders,
			Body:    body,
		}
		return fmt.Sprintf("%s %s%sContent-Type: %s%sContent-Length: %d%s%s%s", response.Version, response.Status, CRLF, response.Headers.ContentType, CRLF, response.Headers.ContentLength, CRLF, CRLF, response.Body)
	}

	if match, _ := path.Match("/echo/*", request.Path); match {
		body := path.Base(request.Path)
		responseHeaders := ResponseHeaders{
			ContentType:   "text/plain",
			ContentLength: len(body),
		}
		response := Response{
			Version: request.Version,
			Status:  fmt.Sprintf("%d %s", http.StatusOK, "OK"),
			Headers: responseHeaders,
			Body:    body,
		}
		return fmt.Sprintf("%s %s%sContent-Type: %s%sContent-Length: %d%s%s%s", response.Version, response.Status, CRLF, response.Headers.ContentType, CRLF, response.Headers.ContentLength, CRLF, CRLF, response.Body)
	}

	if path.Base(request.Path) == "/" {
		response := Response{
			Version: request.Version,
			Status:  fmt.Sprintf("%d %s", http.StatusOK, "OK"),
		}
		return fmt.Sprintf("%s %s%s%s", response.Version, response.Status, CRLF, CRLF)
	}

	response := Response{
		Version: request.Version,
		Status:  fmt.Sprintf("%d %s", http.StatusNotFound, "Not Found"),
	}
	return fmt.Sprintf("%s %s%s%s", response.Version, response.Status, CRLF, CRLF)
}
