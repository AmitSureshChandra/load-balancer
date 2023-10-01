package main

import (
	"log"
	"net"
	"net/http"
	"sync"
)

var backendServers = []string{"localhost:8090", "localhost:8091"}
var curIndex int
var mu sync.Mutex

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		// get backend add in round-robin fashion
		backendAdd := getNextBackend()

		log.Println(backendAdd)

		// create a proxy connection to backend server
		backendConn, err := net.Dial("tcp", backendAdd)

		if err != nil {
			http.Error(writer, "Error connecting to backend", http.StatusServiceUnavailable)
			return
		}
		defer backendConn.Close()

		log.Println("backend connected")

		err = request.Write(backendConn)
		log.Println("req written")

		if err != nil {
			http.Error(writer, "Error Forwarding Req", http.StatusBadGateway)
			return
		}

		resData := make([]byte, 1024)

		_, err = backendConn.Read(resData)
		if err != nil {
			return
		}

		_, err = writer.Write([]byte(backendAdd))

		//_, err = io.Copy(writer, backendConn)
		log.Println("res copied")
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
	if len(backendServers) == 0 {
		return ""
	}
	curIndex := (curIndex + 1) % len(backendServers)
	return backendServers[curIndex]
}
