# accesslogger


[![Documentation](https://godoc.org/github.com/mashiike/accesslogger?status.svg)](https://godoc.org/github.com/mashiike/accesslogger)
![Latest GitHub tag](https://img.shields.io/github/tag/mashiike/accesslogger.svg)
![Github Actions test](https://github.com/mashiike/accesslogger/workflows/Test/badge.svg?branch=main)
[![Go Report Card](https://goreportcard.com/badge/mashiike/accesslogger)](https://goreportcard.com/report/mashiike/accesslogger)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/mashiike/accesslogger/blob/master/LICENSE)

accessloger for golang http handler

## Requirements
  * Go 1.18 or higher. support the 3 latest versions of Go.

## Usage 

sample go code
```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/mashiike/accesslogger"
)

func main() {
	err := http.ListenAndServe("localhost:8080", 
		accesslogger.Wrap(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
				e := json.NewEncoder(w)
				e.SetEscapeHTML(true)
				e.SetIndent("", "  ")
				e.Encode(map[string]string{"status":"ok"})
			}),
			accesslogger.CombinedLogger(&combinedLogs),
			//accesslogger.CombinedDLogger(&combinedDLogs),
			//accesslogger.JSONLogger(&jsonLogs),
        ),
	)
	if err != nil {
		log.Fatalln(err)
	}
}
```

output log:
```
192.0.2.1:1234 - - [26/Dec/2022:15:04:05 +0900] "GET / HTTP/1.1" - 4 "-" "go test client"
192.0.2.1:1234 - hoge [26/Dec/2022:15:04:06 +0900] "GET /hoge HTTP/1.1" - 4 "https://example.com" "go test client"
222.222.333.333 - hoge [26/Dec/2022:15:04:06 +0900] "GET /hoge HTTP/1.1" - 4 "https://example.com" "go test client"
```
