package main

import (
	"gopkg.in/mgo.v2"
	"log"
)

func main() {
	// connect to mongodb
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Panic(err)
	}
	defer session.Close()

	RunServer(session)
}
