package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSimpleAPI(t *testing.T) {
	trackingWords = []string{"obama"}
	ts := httptest.NewServer(GetMainEngine())
	defer ts.Close()
	time.Sleep(5000 * time.Millisecond)
	res, err := http.Get(ts.URL + "/web/obama")
	if err != nil {
		log.Fatal(err)
	}
	_, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
}
