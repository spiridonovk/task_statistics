package main

import (
	"flag"
	"log"
	"net/http"
	"test_task/api"
	"test_task/engine"
)

func main() {
	initServer()
}

func initServer() {
	port := flag.String("port", "8080", "tcp port")
	flag.Parse()
	engine.Init()
	engine.Timer()
	mux := http.NewServeMux()
	mux.HandleFunc("/statistics", api.GetStatistics)
	mux.HandleFunc("/command/", api.SetCommand)

	srv := http.Server{
		Addr:    ":" + *port,
		Handler: mux,
	}
	log.Println("Start http server at : " + *port)
	log.Fatal(srv.ListenAndServe())
}
