package main

import (
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/mvdan/xurls"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type linkData struct {
	sourceTweets []*anaconda.Tweet
	title        string
	URL          string
}

type tweetsDisplay struct {
	tweetsByFav tweetsSlice
	tweetsByRet tweetsSlice
}

type linksDisplay struct {
	linksByFav linksSlice
	linksByRet linksSlice
}

type tweetsMap map[string]*anaconda.Tweet
type linksMap map[string]*linkData
type tweetsSlice []*anaconda.Tweet
type linksSlice []*linkData

type dataChannelValues struct {
	tweets tweetsMap
	links  linksMap
}

type displayData struct {
	tweets *tweetsDisplay
	links  *linksDisplay
}

type wordToTweetMap map[string]tweetsMap
type wordToLinksMap map[string]linksMap

type dump struct {
	tweets wordToTweetMap
	links  wordToLinksMap
}

func getValues(tweetIDtoTweet tweetsMap) []*anaconda.Tweet {
	tweets := make(tweetsSlice, len(tweetIDtoTweet))
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

func expandURLs(urls []string) []*linkData {
	resultingURLs := make([]*linkData, 0)
	for _, url := range urls {
		finalURL := url
		var resp *http.Response
		var err error
		for {
			resp, err = http.Get(finalURL)
			if err == nil {
				if finalURL == resp.Request.URL.String() {
					break
				}
				finalURL = resp.Request.URL.String()
			} else {
				break
			}
		}
		if finalURL != "" {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			r, _ := regexp.Compile("<title>(.*)</title>")
			match := r.FindStringSubmatch(string(body))
			var title string
			if match != nil {
				title = match[1]
			}
			resultingURLs = append(resultingURLs,
				&linkData{sourceTweets: nil, title: title, URL: finalURL})
		} else {
			fmt.Printf("Couldn't find the final URL for %s", url)
		}
	}
	return resultingURLs
}

func storeTweet(tweetMap wordToTweetMap, links wordToLinksMap, tweet *anaconda.Tweet) {
	subWords := make([]string, 0)
	urls := xurls.Relaxed.FindAllString(tweet.Text, -1)
	resultingURLs := expandURLs(urls)
	for _, word := range trackingWords {
		if strings.Contains(strings.ToLower(tweet.Text), word) {
			subWords = append(subWords, word)
			fmt.Printf("Tweet has %s\n", word)
			if tweetMap[word] == nil || links[word] == nil {
				tweetMap[word] = make(map[string]*anaconda.Tweet)
				links[word] = make(map[string]*linkData)
			}
			tweetMap[word][tweet.IdStr] = tweet
			for _, link := range resultingURLs {
				ld, found := links[word][link.URL]
				if found {
					ld.sourceTweets = append(ld.sourceTweets, tweet)
				} else {
					link.sourceTweets = make([]*anaconda.Tweet, 1)
					link.sourceTweets[0] = tweet
					links[word][link.URL] = link
				}
			}
		}
	}
}

func dataManager(req chan<- *dataChannelValues, ask <-chan string) {
	if len(trackingWords) < 1 {
		panic("Need to supply at least one words to track")
	}

	tweets := make(wordToTweetMap)
	links := make(wordToLinksMap)
	dumpContents, err := ioutil.ReadFile("/tweets/dump")
	if err != nil {
		fmt.Println("Failed to restore the dump")
	}
	err = json.Unmarshal(dumpContents, &tweets)
	if err != nil {
		fmt.Println("Failed to unmarshal the dump")
		fmt.Println(err)
	} else {
		fmt.Println("Successfully restored the previous dump")
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
					storeTweet(tweets, links, &t)
				} else {
					// TODO something better for retweets
					originalTweet := t.RetweetedStatus
					storeTweet(tweets, links, originalTweet)
				}
			}
			if tweetsNumber != nil {
				if count <= 0 {
					break
				}
				count = count - 1
			}
		case word := <-ask:
			req <- &dataChannelValues{tweets: tweets[word], links: links[word]}
		case <-dumpReq:
			dumpReq <- &dump{tweets: tweets, links: links}
		}
	}
}
