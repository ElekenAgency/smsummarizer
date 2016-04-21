package main

import (
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/mvdan/xurls"
	"io/ioutil"
	"net/url"
	"strings"
)

type tweetsMap map[string]map[string]*anaconda.Tweet

func getValues(tweetIDtoTweet map[string]*anaconda.Tweet) []*anaconda.Tweet {
	tweets := make([]*anaconda.Tweet, len(tweetIDtoTweet))
	idx := 0
	for key := range tweetIDtoTweet {
		tweets[idx] = tweetIDtoTweet[key]
		idx++
	}
	return tweets
}

func simplifyTweets(tweets []*anaconda.Tweet) []TweetShort {
	result := make([]TweetShort, len(tweets))
	for i, tweet := range tweets {
		result[i] = TweetShort{tweet.Text, tweet.FavoriteCount, tweet.RetweetCount}
	}
	return result
}

func storeTweet(tweetMap tweetsMap, tweet *anaconda.Tweet) {
	for _, word := range trackingWords {
		if strings.Contains(strings.ToLower(tweet.Text), word) {
			fmt.Printf("Tweet has %s\n", word)
			if tweetMap[word] == nil {
				tweetMap[word] = make(map[string]*anaconda.Tweet)
			}
			tweetMap[word][tweet.IdStr] = tweet
		}
	}
	urls := xurls.Relaxed.FindAllString(tweet.Text, -1)
	if len(urls) > 0 && *fullLog {
		fmt.Printf("The tweet has %d URLs\n", len(urls))
	}
}

type TweetsData struct {
	tweetsByFav []*anaconda.Tweet
	tweetsByRet []*anaconda.Tweet
}

func dataManager(req chan<- map[string]*anaconda.Tweet, ask <-chan string) {
	if len(trackingWords) < 1 {
		panic("Need to supply at least one words to track")
	}

	tweets := make(tweetsMap)
	dumpContents, err := ioutil.ReadFile("/tweets/dump")
	if err != nil {
		fmt.Println("Failed to restore the dump")
	}
	err = json.Unmarshal(dumpContents, &tweets)
	if err != nil {
		fmt.Println("Failed to unmarshal the dump")
		fmt.Println(err)
	}
	anaconda.SetConsumerKey("TgFsDmBWfiQb7i0QhyGkgA")
	anaconda.SetConsumerSecret("nDKbC8diEDeYq5ZN4QOv2RhxfyX4UebX0ZtbqPVDU")
	api := anaconda.NewTwitterApi("244167420-jOu3uiiBvZS7m5JkXaDhIQROjc1jooBYgawSD7Q2", "eQHohTUq4e63DlnrxZ9wZ43g7R5eKTX7tau2m0WewjlU2")
	v := url.Values{}
	v.Set("track", strings.Join(trackingWords, ", "))
	fmt.Println("Tracking - " + strings.Join(trackingWords, ", "))
	if *fullLog {
		api.SetLogger(anaconda.BasicLogger)
	}
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
			}
			if tweetsNumber != nil {
				if count <= 0 {
					break
				}
				count = count - 1
			}
		case word := <-ask:
			req <- tweets[word]
		case <-dumpReq:
			dumpRes <- tweets
		}
	}
}
