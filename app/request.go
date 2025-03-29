package main

import (
	"net/http"
	"strings"
)

type RequestHeaders struct {
	Host      string
	UserAgent string
	Accept    string
}

type Request struct {
	Method  string
	Path    string
	Version string
	Headers RequestHeaders
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
		}
	}
	return requestHeaders
}

func DeserializeRequest(request string) Request {
	parts := strings.Split(request, "\r\n")
	status := strings.Split(parts[0], " ")
	headers := parseHeaders(parts[1:])
	var method string
	switch status[0] {
	case "GET":
		method = http.MethodGet
	}
	return Request{
		Method:  method,
		Path:    status[1],
		Version: status[2],
		Headers: headers,
	}
}
