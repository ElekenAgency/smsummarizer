package main

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/bradfitz/slice"
)

func processor(req <-chan interface{}, tweetsToDisplay chan<- *TweetsData) {
	tweetsC := make(chan map[string]*anaconda.Tweet)
	requestData := make(chan interface{})
	go dataManager(tweetsC, requestData)
	for {
		select {
		case tweets := <-tweetsC:
			tweetsToDisplay <- processTweets(tweets)
		case <-req:
			requestData <- 1
		}
	}
}

func processTweets(tweetsMap map[string]*anaconda.Tweet) *TweetsData {
	// we will use map
	// type of max -> array of indexes
	// array of tweets
	// for some reason we need to use the ids because otherwise it appends to the end
	tweets := make([]*anaconda.Tweet, len(tweetsMap))
	id := 0
	for _, tweet := range tweetsMap {
		tweets[id] = tweet
		id++
	}
	// get sorted by favorite
	// somehow there is a memory problem here and I am not sure why
	// try printing first and then using another way of sorting later
	favInd := arrayIndexes(len(tweets))
	slice.Sort(favInd[:], func(i, j int) bool {
		return tweets[i].FavoriteCount < tweets[j].FavoriteCount
	})
	// get sorted by retweets
	retwInd := arrayIndexes(len(tweets))
	slice.Sort(retwInd[:], func(i, j int) bool {
		return tweets[i].RetweetCount < tweets[j].RetweetCount
	})
	return &TweetsData{tweets, favInd, retwInd}
}

func arrayIndexes(len int) []int {
	result := make([]int, len)
	for i := 0; i < len; i++ {
		result[i] = i
	}
	return result
}
