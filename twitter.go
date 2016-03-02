package main

import (
	"flag"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/mvdan/xurls"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
)

var Logger *log.Logger

func storeTweet(tweetMap map[string]*anaconda.Tweet, tweet *anaconda.Tweet) {
	tweetMap[tweet.IdStr] = tweet
	urls := xurls.Relaxed.FindAllString(tweet.Text, -1)
	if len(urls) > 0 {
		Logger.Println("urls - ", urls)
	}
}

func initLog() *log.Logger {
	if *debugingMode {
		return log.New(os.Stdout, "DEBUG:", log.Ldate|log.Ltime|log.Lshortfile)
	}
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}

	return log.New(file,
		"PREFIX: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func cleanup() {
	fmt.Println("\nExiting!")
}

func init() {
	flag.Var(&trackingWords, "words", "Words to track")
	// setup listening to CTRL-C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()
}

func main() {
	flag.Parse()
	Logger = initLog()
	if len(trackingWords) < 1 {
		panic("Need to supply at least one words to track")
	}
	tweets := make(map[string]*anaconda.Tweet)
	anaconda.SetConsumerKey("TgFsDmBWfiQb7i0QhyGkgA")
	anaconda.SetConsumerSecret("nDKbC8diEDeYq5ZN4QOv2RhxfyX4UebX0ZtbqPVDU")
	api := anaconda.NewTwitterApi("244167420-jOu3uiiBvZS7m5JkXaDhIQROjc1jooBYgawSD7Q2", "eQHohTUq4e63DlnrxZ9wZ43g7R5eKTX7tau2m0WewjlU2")
	v := url.Values{}
	v.Set("track", strings.Join(trackingWords, ", "))
	Logger.Println("Tracking - " + strings.Join(trackingWords, ", "))
	stream := api.PublicStreamFilter(v)

	count := 0
	if tweetsNumber != nil {
		count = *tweetsNumber
	}
	for o := range stream.C {
		t, ok := o.(anaconda.Tweet)
		if ok {
			if t.RetweetedStatus == nil {
				Logger.Print("Original:\t")
				storeTweet(tweets, &t)
			} else {
				Logger.Print("Retweet:\t")
				originalTweet := t.RetweetedStatus
				Logger.Print("Retweet count:", originalTweet.RetweetCount, "\n")
				sourceTweet, ok := tweets[originalTweet.IdStr]
				if ok {
					sourceTweet.RetweetCount = originalTweet.RetweetCount
				} else {
					storeTweet(tweets, originalTweet)
				}
			}
			Logger.Println(t.Text)
		}
		if tweetsNumber != nil {
			if count <= 0 {
				break
			}
			count = count - 1
		}
	}
}
