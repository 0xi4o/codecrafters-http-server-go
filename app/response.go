package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
	"os"
	"path"
)

type ResponseHeaders struct {
	ContentEncoding string
	ContentLength   int
	ContentType     string
}

type Response struct {
	Version        string
	Status         string
	Headers        ResponseHeaders
	Body           string
	CompressedBody []byte
}

func compressBody(body string) []byte {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	zw.Write([]byte(body))
	zw.Close()
	return buf.Bytes()
}

func SerializeResponse(request Request, flags map[string]string) string {
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
		if request.Headers.AcceptEncoding == "gzip" {
			response.CompressedBody = compressBody(body)
			response.Headers.ContentEncoding = request.Headers.AcceptEncoding
			response.Headers.ContentLength = len(response.CompressedBody)
			return fmt.Sprintf("%s %s%sContent-Encoding: %s%sContent-Type: %s%sContent-Length: %d%s%s%s", response.Version, response.Status, CRLF, response.Headers.ContentEncoding, CRLF, response.Headers.ContentType, CRLF, response.Headers.ContentLength, CRLF, CRLF, response.CompressedBody)
		}
		return fmt.Sprintf("%s %s%sContent-Type: %s%sContent-Length: %d%s%s%s", response.Version, response.Status, CRLF, response.Headers.ContentType, CRLF, response.Headers.ContentLength, CRLF, CRLF, response.Body)
	}

	if match, _ := path.Match("/files/*", request.Path); match {
		switch request.Method {
		case http.MethodGet:
			filename := path.Base(request.Path)
			filepath := fmt.Sprintf("%s%s", flags["directory"], filename)
			fmt.Println(filepath)
			data, err := os.ReadFile(filepath)
			if err != nil {
				response := Response{
					Version: request.Version,
					Status:  fmt.Sprintf("%d %s", http.StatusNotFound, "Not Found"),
				}
				return fmt.Sprintf("%s %s%s%s", response.Version, response.Status, CRLF, CRLF)
			}
			body := string(data)
			responseHeaders := ResponseHeaders{
				ContentEncoding: request.Headers.AcceptEncoding,
				ContentType:     "application/octet-stream",
				ContentLength:   len(body),
			}
			response := Response{
				Version: request.Version,
				Status:  fmt.Sprintf("%d %s", http.StatusOK, "OK"),
				Headers: responseHeaders,
				Body:    body,
			}
			if request.Headers.AcceptEncoding == "gzip" {
				response.CompressedBody = compressBody(body)
				response.Headers.ContentEncoding = request.Headers.AcceptEncoding
				response.Headers.ContentLength = len(response.CompressedBody)
				return fmt.Sprintf("%s %s%sContent-Encoding: %s%sContent-Type: %s%sContent-Length: %d%s%s%s", response.Version, response.Status, CRLF, response.Headers.ContentEncoding, CRLF, response.Headers.ContentType, CRLF, response.Headers.ContentLength, CRLF, CRLF, response.CompressedBody)
			}
			return fmt.Sprintf("%s %s%sContent-Type: %s%sContent-Length: %d%s%s%s", response.Version, response.Status, CRLF, response.Headers.ContentType, CRLF, response.Headers.ContentLength, CRLF, CRLF, response.Body)
		case http.MethodPost:
			response := Response{
				Version: request.Version,
				Status:  fmt.Sprintf("%d %s", http.StatusCreated, "Created"),
			}
			return fmt.Sprintf("%s %s%s%s", response.Version, response.Status, CRLF, CRLF)
		}
	}

	if match, _ := path.Match("/user-agent", request.Path); match {
		body := request.Headers.UserAgent
		responseHeaders := ResponseHeaders{
			ContentEncoding: request.Headers.AcceptEncoding,
			ContentType:     "text/plain",
			ContentLength:   len(body),
		}
		response := Response{
			Version: request.Version,
			Status:  fmt.Sprintf("%d %s", http.StatusOK, "OK"),
			Headers: responseHeaders,
			Body:    body,
		}
		if request.Headers.AcceptEncoding == "gzip" {
			response.CompressedBody = compressBody(body)
			response.Headers.ContentEncoding = request.Headers.AcceptEncoding
			response.Headers.ContentLength = len(response.CompressedBody)
			return fmt.Sprintf("%s %s%sContent-Encoding: %s%sContent-Type: %s%sContent-Length: %d%s%s%s", response.Version, response.Status, CRLF, response.Headers.ContentEncoding, CRLF, response.Headers.ContentType, CRLF, response.Headers.ContentLength, CRLF, CRLF, response.CompressedBody)
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
