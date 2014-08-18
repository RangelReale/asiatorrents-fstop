package main

import (
	"flag"
	"fmt"
	"github.com/RangelReale/asiatorrents-fstop"
	"github.com/RangelReale/filesharetop/importer"
	"gopkg.in/mgo.v2"
	"log"
	"os"
)

var version = flag.Bool("version", false, "show version and exit")

func main() {
	flag.Parse()

	// output version
	if *version {
		fmt.Printf("asiatorrents-importer version %s\n", fstopimp.VERSION)
		os.Exit(0)
	}

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
	imp.Database = "fstop_asiatorrents"

	// create fetcher
	fetcher := asiatorrents.NewFetcher()

	// import data
	err = imp.Import(fetcher)
	if err != nil {
		logger.Fatal(err)
	}

	// consolidate data
	err = imp.Consolidate("", 48)
	if err != nil {
		logger.Fatal(err)
	}

	err = imp.Consolidate("weekly", 168)
	if err != nil {
		logger.Fatal(err)
	}
}
