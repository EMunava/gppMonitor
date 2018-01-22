package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type routerObject struct {
	Mux *mux.Router
}

var router *routerObject

func init() {
	router = &routerObject{}
	go func() {
		log.Println("Starting HTTP Server...")
		log.Fatal(http.ListenAndServe(":8080", getRouter()))
	}()
	defer func() {
		if err := recover(); err != nil {
			fmt.Print(err)
		}
	}()
}

func getRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/dateRollover", runSelenium)
	r.HandleFunc("/retrieveEDOLog", retrieveEDOLog)
	r.HandleFunc("/listFiles", listFiles)

	router.Mux = r
	return r
}
