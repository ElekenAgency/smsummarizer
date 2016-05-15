package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type dReqChan chan interface{}
type dRespChan chan *dump

var dReqC dReqChan
var dRespC dRespChan
var tweetsDumpPath string
var linksDumpPath string

type dump struct {
	tweets wordToTweetMap
	links  wordToLinksMap
}

func readDumpContents(path string, holder interface{}) {
	if _, err := os.Stat(path); err == nil {
		dumpContents, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Printf("Failed to restore the dump from %s\n", path)
		}
		err = json.Unmarshal(dumpContents, &holder)
		if err != nil {
			fmt.Println("Failed to unmarshal the dump")
			fmt.Println(err)
		} else {
			fmt.Println("Successfully restored the previous dump")
		}
	} else {
		fmt.Printf("%s doesn't exists, can't restore dump\n", path)
	}
}

func writeDumpContents(path string, holder interface{}) {
	jsonVal, err := json.Marshal(holder)
	err = ioutil.WriteFile(path, jsonVal, 0644)
	if err != nil {
		fmt.Printf("Problems with saving the data at %s\n", path)
		fmt.Println(err)
	} else {
		fmt.Println("Successfully wrote the dump")
	}
}
