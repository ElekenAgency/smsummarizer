package main

import (
	"github.com/ChimeraCoder/anaconda"
	"sort"
)

func processor(req <-chan string, displayChannel chan<- *displayData) {
	dataChannel := make(chan *dataChannelValues)
	requestData := make(chan string)
	go dataManager(dataChannel, requestData)
	for {
		select {
		case tweetsAndLinks := <-dataChannel:
			displayChannel <- &displayData{tweets: processTweets(tweetsAndLinks.tweets),
				links: processLinks(tweetsAndLinks.links)}
		case word := <-req:
			requestData <- word
		}
	}
}

func processLinks(lm linksMap) *linksDisplay {
	return nil
}

type ByFav []*anaconda.Tweet

func (a ByFav) Len() int           { return len(a) }
func (a ByFav) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFav) Less(i, j int) bool { return a[i].FavoriteCount > a[j].FavoriteCount }

type ByRet []*anaconda.Tweet

func (a ByRet) Len() int           { return len(a) }
func (a ByRet) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRet) Less(i, j int) bool { return a[i].RetweetCount > a[j].RetweetCount }

func processTweets(tweetsMap map[string]*anaconda.Tweet) *tweetsDisplay {
	// we will use map
	// type of max -> array of indexes
	// array of tweets
	// for some reason we need to use the ids because otherwise it appends to the end
	tweets := getValues(tweetsMap)
	tweetsByFav := make([]*anaconda.Tweet, len(tweetsMap))
	tweetsByRet := make([]*anaconda.Tweet, len(tweetsMap))
	copy(tweetsByFav, tweets)
	copy(tweetsByRet, tweets)
	sort.Sort(ByFav(tweetsByFav))
	sort.Sort(ByRet(tweetsByRet))
	return &tweetsDisplay{tweetsByFav, tweetsByRet}
}

func arrayIndexes(len int) []int {
	result := make([]int, len)
	for i := 0; i < len; i++ {
		result[i] = i
	}
	return result
}
