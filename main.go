package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var backendServers = []string{"http://localhost:9000", "http://localhost:9001"}
var curIndex int
var mu sync.Mutex

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		// get backend add in round-robin fashion
		backendAdd := getNextBackend()

		log.Println(backendAdd)

		// create a proxy connection to backend server
		backendResp, err := http.Get(backendAdd)

		if err != nil {
			http.Error(writer, "Error connecting to backend", http.StatusServiceUnavailable)
			return
		}

		defer backendResp.Body.Close()

		_, err = io.Copy(writer, backendResp.Body)

		if err != nil {
			return
		}
	})

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		return
	}
}

func getNextBackend() string {
	mu.Lock()
	defer mu.Unlock()
	n := len(backendServers)
	if n == 0 {
		return ""
	}
	curIndex++
	curIndex %= n
	log.Println("curIndex " + strconv.Itoa(curIndex))
	return backendServers[curIndex]
}
