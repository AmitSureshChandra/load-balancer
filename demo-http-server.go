package main

import (
	"log"
	"net/http"
	"os"
)

func main() {

	port := os.Args[1]

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		log.Println("got req " + port)
		writer.Write([]byte(port + ""))
	})

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		return
	}
}
