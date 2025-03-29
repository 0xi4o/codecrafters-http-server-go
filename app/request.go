package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

type RequestHeaders struct {
	Host          string
	UserAgent     string
	Accept        string
	ContentLength int
	ContentType   string
}

type Request struct {
	Method  string
	Path    string
	Version string
	Headers RequestHeaders
	Body    string
}

func parseHeaders(headers []string) RequestHeaders {
	requestHeaders := RequestHeaders{}
	for _, header := range headers {
		parts := strings.Split(header, ": ")
		switch parts[0] {
		case "Host":
			requestHeaders.Host = parts[1]
		case "User-Agent":
			requestHeaders.UserAgent = parts[1]
		case "Accept":
			requestHeaders.Accept = parts[1]
		case "Content-Length":
			requestHeaders.ContentLength, _ = strconv.Atoi(parts[1])
		case "Content-Type":
			requestHeaders.ContentType = parts[1]
		}
	}
	return requestHeaders
}

func DeserializeRequest(request string) Request {
	parts := strings.Split(request, fmt.Sprintf("%s%s", CRLF, CRLF))
	headerParts := strings.Split(parts[0], CRLF)
	status := strings.Split(headerParts[0], " ")
	headers := parseHeaders(headerParts[1:])
	var method string
	switch status[0] {
	case "GET":
		method = http.MethodGet
	case "POST":
		method = http.MethodPost
	}
	return Request{
		Method:  method,
		Path:    status[1],
		Version: status[2],
		Headers: headers,
		Body:    parts[1],
	}
}

func (request *Request) Process(flags map[string]string) error {
	if match, _ := path.Match("/files/*", request.Path); match {
		if request.Method == http.MethodPost {
			filename := path.Base(request.Path)
			filepath := fmt.Sprintf("%s%s", flags["directory"], filename)
			fmt.Println(filepath)
			err := os.WriteFile(filepath, []byte(request.Body), 0666)
			if err != nil {
				return err
			}
			err = os.Truncate(filepath, int64(request.Headers.ContentLength))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
