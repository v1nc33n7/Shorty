package main

import (
	"log"
)

func main() {
	err := ConnRedis(":6379")
	if err != nil {
		log.Fatal(err)
	}

	RunServer()
}
