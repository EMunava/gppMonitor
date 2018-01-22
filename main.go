package main

import (
	_ "github.com/CardFrontendDevopsTeam/GPPMonitor/rest"
	//_ "github.com/CardFrontendDevopsTeam/GPPMonitor/selenium"
	"log"
	"time"
)

func main() {
	log.Println("GPP Monitor")
	for true {
		time.Sleep(10 * time.Minute)
	}
}
