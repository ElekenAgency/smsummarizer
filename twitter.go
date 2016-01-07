package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/mvdan/xurls"
	"net/url"
	"os"
	"strings"
)

func main() {
	trackWords := os.Args[1:]
	if len(trackWords) < 1 {
		panic("Need to supply at least one words to track")
	}
	tweets := make(map[string]anaconda.Tweet)
	anaconda.SetConsumerKey("TgFsDmBWfiQb7i0QhyGkgA")
	anaconda.SetConsumerSecret("nDKbC8diEDeYq5ZN4QOv2RhxfyX4UebX0ZtbqPVDU")
	api := anaconda.NewTwitterApi("244167420-jOu3uiiBvZS7m5JkXaDhIQROjc1jooBYgawSD7Q2", "eQHohTUq4e63DlnrxZ9wZ43g7R5eKTX7tau2m0WewjlU2")
	v := url.Values{}
	v.Set("track", strings.Join(trackWords, ", "))
	fmt.Println("Tracking - " + strings.Join(trackWords, ", "))
	api.SetLogger(anaconda.BasicLogger)
	stream := api.PublicStreamFilter(v)
	k := 0
	for o := range stream.C {
		t, ok := o.(anaconda.Tweet) // try casting into a tweet
		if ok {
			k++
			if !t.Retweeted {
				fmt.Print("Original:\t")
				tweets[t.IdStr] = t
				urls := xurls.Relaxed.FindAllString(t.Text, -1)
				fmt.Print("urls - ")
				fmt.Println(urls)
			} else {
				fmt.Print("Retweet:\t")
				sourceTweet, ok := tweets[t.Source]
				if ok {
					sourceTweet.RetweetCount++
				}
			}
			fmt.Println(t.Text)
		}
		if k > 10 {
			break
		}
	}
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	for key, value := range tweets {
		fmt.Println("Key:", key, "Value:", value.Text, "IsRetweet", value.Retweeted)
	}
	// v := url.Values{}
	// v.Set("count", "30")
	// searchResult, _ := api.GetSearch("golang", v)
	// for _, tweet := range searchResult.Statuses {
	// 	fmt.Println(tweet.User.Name + ":" + tweet.Text)
	// 	urls := xurls.Relaxed.FindAllString(tweet.Text, -1)
	// 	fmt.Print("urls - ")
	// 	fmt.Println(urls)
	// 	fmt.Println(fmt.Sprintf("favorited - %d, retweeted - %d", tweet.FavoriteCount, tweet.RetweetCount))
	// 	if tweet.RetweetedStatus != nil {
	// 		fmt.Println("Retweeted!!!")
	// 	} else {
	// 		fmt.Println("Original")
	// 	}
	// }
}
