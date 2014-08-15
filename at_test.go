package asiatorrents

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestFetcher(t *testing.T) {
	f := NewFetcher()
	f.SetLogger(log.New(ioutil.Discard, "", 0))
	i, err := f.Fetch()
	if err != nil {
		t.Error(err)
		return
	}

	if len(i) == 0 {
		t.Error("No data returned from parser")
	}
}
