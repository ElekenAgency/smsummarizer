package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/mvdan/xurls"
	"log"
	"net/url"
	"strings"
)

func storeTweet(tweetMap map[string]map[string]*anaconda.Tweet, tweet *anaconda.Tweet) {
	for _, word := range trackingWords {
		if strings.Contains(strings.ToLower(tweet.Text), word) {
			if tweetMap[word] == nil {
				tweetMap[word] = make(map[string]*anaconda.Tweet)
			}
			tweetMap[word][tweet.IdStr] = tweet
		}
	}
	urls := xurls.Relaxed.FindAllString(tweet.Text, -1)
	if len(urls) > 0 {
		Logger.Println("urls - ", urls)
	}
}

func summary(tweets map[string]map[string]*anaconda.Tweet) {
	Logger.Print("=== Summary ===")
	for word := range tweets {
		Logger.Print("\t" + word)
		for _, tweet := range tweets[word] {
			Logger.Printf("\t\t%s\n", tweet.Text)
		}
	}
}

func dataManager(words []string, log *log.Logger, tweetsNumber *int, req chan<- *anaconda.Tweet, ask <-chan int) {
	trackingWords = words
	Logger = log
	if len(trackingWords) < 1 {
		panic("Need to supply at least one words to track")
	}

	tweets := make(map[string]map[string]*anaconda.Tweet)
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
	for {
		select {
		case o := <-stream.C:
			t, ok := o.(anaconda.Tweet)
			if ok {
				if t.RetweetedStatus == nil {
					storeTweet(tweets, &t)
				} else {
					// TODO something better for retweets
					originalTweet := t.RetweetedStatus
					storeTweet(tweets, originalTweet)
				}
				Logger.Println(t.Text)
			}
			if tweetsNumber != nil {
				if count <= 0 {
					break
				}
				count = count - 1
			}
		case <-ask:
			fmt.Println("asking for a tweet ")
			for _, val := range tweets[trackingWords[0]] {
				req <- val
				break
			}
		}
	}
}
