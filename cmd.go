package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

type words []string

// String is the method to format the flag's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (w *words) String() string {
	return fmt.Sprint(*w)
}

// Set is the method to set the flag value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the flag.
// It's a comma-separated list, so we split it.
func (w *words) Set(value string) error {
	if len(*w) > 0 {
		return errors.New("words flag already set")
	}
	for _, dt := range strings.Split(value, ",") {
		trackingWord := dt
		*w = append(*w, trackingWord)
	}
	return nil
}

var debugingMode = flag.Bool("debug", false, "Debugging mode")
var testing = flag.Bool("testing", false, "Testing mode")
var tweetsNumber = flag.Int("tweets", -1, "Number of tweets to be captured. By default, infinite")
var trackingWords words
