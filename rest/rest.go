package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type RouterObject struct {
	Mux *mux.Router
}

var router *RouterObject

func init() {
	router = &RouterObject{}
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
	r.HandleFunc("/dateRollover", runSelenium).Methods("POST")

	router.Mux = r
	return r
}

/*
Router starts the router service
*/
func Router() *RouterObject {
	return router
}
