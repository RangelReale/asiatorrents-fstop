package main

import (
	"github.com/RangelReale/asiatorrents-fstop"
	"github.com/RangelReale/filesharetop/importer"
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

	// create and run importer
	imp := fstopimp.NewImporter(logger, session)

	// create fetcher
	fetcher := asiatorrents.NewFetcher()

	// import data
	err = imp.Import(fetcher)
	if err != nil {
		logger.Fatal(err)
	}

	// consolidate data
	err = imp.Consolidate(48)
	if err != nil {
		logger.Fatal(err)
	}
}
