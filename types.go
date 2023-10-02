package main

import (
	"log"
	"strconv"
	"sync"
)

type Config struct {
	RouteStrategyCode string `json:"route_strategy_code"`
}

type Server struct {
	url string
}

type RouteStrategy interface {
	GetBackendServer() *Server
}

type RoundRobinStrategy struct {
	curIndex int
	mu       sync.Mutex
}

func (roundRobin *RoundRobinStrategy) GetBackendServer() *Server {
	roundRobin.mu.Lock()
	defer roundRobin.mu.Unlock()
	n := len(backendServers)
	if n == 0 {
		return &Server{}
	}
	roundRobin.curIndex++
	roundRobin.curIndex %= n
	log.Println("roundRobin.curIndex " + strconv.Itoa(roundRobin.curIndex))
	return &backendServers[roundRobin.curIndex]
}
