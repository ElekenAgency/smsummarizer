package main

import (
	"github.com/ChimeraCoder/anaconda"
	"sort"
)

func processor(req <-chan string, tweetsToDisplay chan<- *TweetsData) {
	tweetsC := make(chan map[string]*anaconda.Tweet)
	requestData := make(chan string)
	go dataManager(tweetsC, requestData)
	for {
		select {
		case tweets := <-tweetsC:
			tweetsToDisplay <- processTweets(tweets)
		case word := <-req:
			requestData <- word
		}
	}
}

type ByFav []*anaconda.Tweet

func (a ByFav) Len() int           { return len(a) }
func (a ByFav) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFav) Less(i, j int) bool { return a[i].FavoriteCount > a[j].FavoriteCount }

type ByRet []*anaconda.Tweet

func (a ByRet) Len() int           { return len(a) }
func (a ByRet) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRet) Less(i, j int) bool { return a[i].RetweetCount > a[j].RetweetCount }

func processTweets(tweetsMap map[string]*anaconda.Tweet) *TweetsData {
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
	return &TweetsData{tweetsByFav, tweetsByRet}
}

func arrayIndexes(len int) []int {
	result := make([]int, len)
	for i := 0; i < len; i++ {
		result[i] = i
	}
	return result
}
