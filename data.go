package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/mvdan/xurls"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

var mutexLinks = &sync.Mutex{}
var mutexTweets = &sync.Mutex{}

type linkData struct {
	sourceTweets []*anaconda.Tweet
	Title        string
	URL          string
	Retweets     int
	Likes        int
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

func (lm linksMap) Put(key string, value *linkData) {
	mutexLinks.Lock()
	lm[key] = value
	mutexLinks.Unlock()
}

func (lm linksMap) Get(key string) *linkData {
	mutexLinks.Lock()
	val := lm[key]
	mutexLinks.Unlock()
	return val
}

func (tm tweetsMap) Put(key string, value *anaconda.Tweet) {
	mutexTweets.Lock()
	tm[key] = value
	mutexTweets.Unlock()
}

func (tm tweetsMap) Get(key string) *anaconda.Tweet {
	mutexTweets.Lock()
	val := tm[key]
	mutexTweets.Unlock()
	return val
}

func getLinksValues(lm linksMap) linksSlice {
	links := make(linksSlice, len(lm))
	idx := 0
	for key := range lm {
		links[idx] = lm.Get(key)
		idx++
	}
	return links
}

func getTweetValues(tweetIDtoTweet tweetsMap) tweetsSlice {
	tweets := make(tweetsSlice, len(tweetIDtoTweet))
	idx := 0
	for key := range tweetIDtoTweet {
		tweets[idx] = tweetIDtoTweet.Get(key)
		idx++
	}
	return tweets
}

func expandURLs(urls []string) linksSlice {
	resultingURLs := make(linksSlice, 0)
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
		if finalURL != "" && resp != nil {
			body, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			r, _ := regexp.Compile(".*<title>(.*)</title>.*")
			match := r.FindStringSubmatch(string(body))
			title := "No title"
			if len(match) > 1 {
				title = match[1]
			}
			resultingURLs = append(resultingURLs,
				&linkData{sourceTweets: nil, Title: title, URL: finalURL})
		} else {
			fmt.Printf("Couldn't find the final URL for %s", url)
		}
	}
	return resultingURLs
}

func contains(s tweetsSlice, e *anaconda.Tweet) (*anaconda.Tweet, int) {
	for i, a := range s {
		if a.IdStr == e.IdStr {
			return a, i
		}
	}
	return nil, 0
}

func storeTweet(tweetMap wordToTweetMap, links wordToLinksMap, tweet *anaconda.Tweet) {
	subWords := make([]string, 0)
	urls := xurls.Relaxed.FindAllString(tweet.Text, -1)
	resultingURLs := expandURLs(urls)
	for _, word := range trackedWords {
		if strings.Contains(strings.ToLower(tweet.Text), word) {
			subWords = append(subWords, word)
			if tweetMap[word] == nil || links[word] == nil {
				tweetMap[word] = make(tweetsMap)
				links[word] = make(linksMap)
			}
			tweetMap[word].Put(tweet.IdStr, tweet)
			for _, link := range resultingURLs {
				ld, found := links[word][link.URL]
				if found {
					// check if it is already there
					elem, id := contains(ld.sourceTweets, tweet)
					if elem != nil {
						ld.sourceTweets = append(ld.sourceTweets[:id], ld.sourceTweets[id+1:]...)
						ld.sourceTweets = append(ld.sourceTweets, tweet)
						ld.Retweets = ld.Retweets - elem.RetweetCount + tweet.RetweetCount
						ld.Likes = ld.Likes - elem.RetweetCount + tweet.FavoriteCount
					} else {
						ld.sourceTweets = append(ld.sourceTweets, tweet)
						ld.Retweets = ld.Retweets + tweet.RetweetCount
						ld.Likes = ld.Likes + tweet.FavoriteCount
					}
					links[word].Put(link.URL, ld)
				} else {
					link.sourceTweets = make(tweetsSlice, 1)
					link.sourceTweets[0] = tweet
					link.Retweets = tweet.RetweetCount
					link.Likes = tweet.FavoriteCount
					links[word].Put(link.URL, link)
				}
			}
		}
	}
}

func dataManager(pRespC chan<- *dataChannelValues, pReqC <-chan string) {
	if len(trackedWords) < 1 {
		panic("Need to supply at least one words to track")
	}
	// create structures and restore previous state
	wtm := make(wordToTweetMap)
	wlm := make(wordToLinksMap)
	readDumpContents(tweetsDumpPath, wtm)
	readDumpContents(linksDumpPath, wlm)
	// set Twitter credentials
	anaconda.SetConsumerKey("TgFsDmBWfiQb7i0QhyGkgA")
	anaconda.SetConsumerSecret("nDKbC8diEDeYq5ZN4QOv2RhxfyX4UebX0ZtbqPVDU")
	api := anaconda.NewTwitterApi("244167420-jOu3uiiBvZS7m5JkXaDhIQROjc1jooBYgawSD7Q2", "eQHohTUq4e63DlnrxZ9wZ43g7R5eKTX7tau2m0WewjlU2")
	// set tracking parameters
	v := url.Values{}
	v.Set("track", strings.Join(trackedWords, ", "))
	stream := api.PublicStreamFilter(v)
	// loop to process requests
	for {
		select {
		case <-dReqC:
			stream.Stop()
			dRespC <- &dump{tweets: wtm, links: wlm}
			return
		case word := <-pReqC:
			pRespC <- &dataChannelValues{tweets: wtm[word], links: wlm[word]}
		case o := <-stream.C:
			t, ok := o.(anaconda.Tweet)
			if ok {
				if t.RetweetedStatus == nil {
					go storeTweet(wtm, wlm, &t)
				} else {
					originalTweet := t.RetweetedStatus
					go storeTweet(wtm, wlm, originalTweet)
				}
			}
		default:
		}
	}
}
