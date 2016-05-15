package main

import (
	"fmt"
)

func addChannel(client *Client, data interface{}) {
	var channel Channel
	mapstructure.Decode(data, &channel)
}
