package main

import (
	"github.com/RangelReale/filesharetop/site"
	"gopkg.in/mgo.v2"
	"log"
	"os"
)

func main() {
	// connect to mongodb
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Panic(err)
	}
	defer session.Close()

	// create logger
	logger := log.New(os.Stderr, "", log.LstdFlags)

	config := fstopsite.NewConfig(13111)
	config.Title = "AsiaTorrents Top"
	config.Logger = logger
	config.Session = session
	config.Database = "fstop_asiatorrents"
	config.TopId = "weekly"

	fstopsite.RunServer(config)
}
