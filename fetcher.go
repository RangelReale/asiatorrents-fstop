package asiatorrents

import (
	"github.com/RangelReale/filesharetop/lib"
	"io/ioutil"
	"log"
)

type Fetcher struct {
	logger *log.Logger
	config *Config
}

func NewFetcher(config *Config) *Fetcher {
	return &Fetcher{
		logger: log.New(ioutil.Discard, "", 0),
		config: config,
	}
}

func (f *Fetcher) ID() string {
	return "ASIATORRENTS"
}

func (f *Fetcher) SetLogger(l *log.Logger) {
	f.logger = l
}

func (f *Fetcher) Fetch() (map[string]*fstoplib.Item, error) {
	parser := NewATParser(f.config, f.logger)

	// parse 4 pages ordered by seeders
	err := parser.Parse(ATSORT_SEEDERS, ATSORTBY_DESCENDING, 4)
	if err != nil {
		return nil, err
	}

	// parse 2 pages ordered by leechers
	err = parser.Parse(ATSORT_LEECHERS, ATSORTBY_DESCENDING, 2)
	if err != nil {
		return nil, err
	}

	return parser.List, nil
}

func (f *Fetcher) CategoryMap() (*fstoplib.CategoryMap, error) {
	return &fstoplib.CategoryMap{
		"MOVIE": []string{"9", "1", "3", "8", "2", "11", "16", "15", "5", "6", "70", "71", "72", "73", "120", "74", "75", "76", "77"},
		"TV":    []string{"20", "19", "21", "22", "23", "28", "29", "79", "80", "81", "82", "83", "84", "85"},
		"MUSIC": []string{"24", "25", "26", "68", "119", "27"},
		"CAT3":  []string{"109", "101", "108", "102", "107", "106", "105", "104", "103"},
	}, nil
}
