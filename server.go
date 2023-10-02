package main

import (
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"net/http"
	"os"
)

var config Config

var routeStrategy RouteStrategy

var defaultStrategy = "round_robin"

var backendServers = []Server{
	{url: "http://localhost:9000"},
	{url: "http://localhost:9001"},
}

func setUpHttpServer() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		// get backend add in round-robin fashion
		backendServer := routeStrategy.GetBackendServer()

		log.Println(backendServer)

		// create a proxy connection to backend server
		backendResp, _ := http.Get(backendServer.url)
		defer backendResp.Body.Close()
		_, _ = io.Copy(writer, backendResp.Body)
	})

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		return
	}
}

func setUpRouteStrategy() {
	file, err := os.Open("config.yaml")
	if err != nil {
		log.Fatal("failed to read config.yaml")
	}

	defer file.Close()

	err = yaml.NewDecoder(file).Decode(&config)
	if err != nil {
		log.Fatal("failed to decode config file")
	}

	if config.RouteStrategyCode == "" {
		config.RouteStrategyCode = defaultStrategy
	}

	switch config.RouteStrategyCode {
	case "round_robin":
		routeStrategy = &RoundRobinStrategy{}
	}
}
