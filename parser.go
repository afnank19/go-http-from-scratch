package main

import (
	"errors"
	"log"
	"strings"
)

// GET /hello HTTP/1.1\r\n
// Host: example.com\r\n
// User-Agent: curl/7.88.1\r\n
// Accept: */*\r\n
// \r\n

func ParseHTTPRequest(reqStr string) *Request {
	req := strings.Split(reqStr, "\r\n")

	var headers []string
	if len(req) > 2 {
		headers = req[:len(req)-2]
	}

	requestLine := headers[0] // first line will always be a request line

	method, path, httpVersion, err := parseRequestLine(requestLine)
	if err != nil {
		log.Fatalln(err.Error()) // handle better
	}

	hdrs := parseHeaders(headers)

	if httpVersion == "HTTP/1.1" && !hdrs.hostHeaderExists {
		log.Fatalln("400 Bad Request")
	}

	// Can expand the request to include the headers/body as well later on
	return &Request{
		method:      method,
		path:        path,
		httpVersion: httpVersion,
		headers:     *hdrs,
	}

}

func parseRequestLine(requestLine string) (string, string, string, error) {
	reqLineTokens := strings.Split(requestLine, " ")

	method := reqLineTokens[0] // GET, POST etc
	path := reqLineTokens[1]
	httpVersion := reqLineTokens[2]

	// Currently only HTTP/1.1
	if httpVersion != "HTTP/1.1" {
		return "", "", "", errors.New("HTTP version mismatch/invalid")
	}

	return method, path, httpVersion, nil
}

func parseHeaders(headers []string) *Headers {
	var hostHeaderExists bool = false
	var acceptType string
	var userAgent string
	for i, r := range headers {
		// fmt.Println("Parsing -> ", r)
		if i == 0 {
			continue
		}

		headerTokens := strings.Split(r, " ")

		if headerTokens[0] == "Host:" {
			hostHeaderExists = true
		}

		if headerTokens[0] == "User-Agent:" {
			userAgent = headerTokens[1]
		}

		if headerTokens[0] == "Accept:" {
			acceptType = headerTokens[1]
		}
	}

	return &Headers{
		hostHeaderExists: hostHeaderExists,
		acceptType:       acceptType,
		userAgent:        userAgent,
	}
}
