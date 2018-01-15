package main

import (
	"log"
	_ "github.com/CardFrontendDevopsTeam/GPPMonitor/selenium"
	"time"
)

func main() {
	log.Println("GPP Monitor")
	for true {
		time.Sleep(10 * time.Minute)
	}
}
