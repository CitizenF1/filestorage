package main

import (
	"log"
	metabus "metabus/src"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	storage := metabus.CreateStorage()
	addresseble := metabus.CreateAddresseble()
	tcm := metabus.CreateTCM()

	r := mux.NewRouter()

	// http storage methods
	r.HandleFunc("/storage", storage.FuncHandler)
	r.HandleFunc("/tcm", tcm.UploadHandler)

	// http serve directory
	r.HandleFunc("/storage/{path:.*}", storage.ServeFiles)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./fm")))

	r.HandleFunc("/StandaloneWindows/{path:.*}", addresseble.ServeFilesAdress)

	server := &http.Server{
		Addr:         ":7360",
		Handler:      r,
		IdleTimeout:  time.Second * 10,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	log.Fatal(server.ListenAndServe())
}
