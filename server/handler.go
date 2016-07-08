package main

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type TrackedWord struct {
	tagName string
}

func listTag(client *Client, data interface{}) {
	fmt.Println("list tag")
	client.send <- Message{"tag list", trackedWords}
}

func addTag(client *Client, data interface{}) {
	var tw TrackedWord
	mapstructure.Decode(data, &tw)
	// TODO add new tags and restart the tracking
}

func updateTag(client *Client, data interface{}) {
	var tw TrackedWord
	if err := mapstructure.Decode(data, &tw); err != nil {
		fmt.Println(err)
	}
	// TODO fix using mapstructure
	dc := data.(map[string]interface{})
	word := dc["tagName"].(string)
	fmt.Println(word)
	dataToDisplay := getDispayData(word, respC, reqC)
	// HERE proper JSON encoding is needed
	client.send <- Message{"tag update", dataToDisplay}
}

func unsubscribeTag(client *Client, data interface{}) {
	var tw TrackedWord
	mapstructure.Decode(data, &tw)
	// TODO add later when we have push updates
}
