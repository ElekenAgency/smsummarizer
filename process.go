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
		case <-dReqC:
			return
		case tweetsAndLinks := <-dataChannel:
			displayChannel <- &displayData{tweets: processTweets(tweetsAndLinks.tweets),
				links: processLinks(tweetsAndLinks.links)}
		case word := <-req:
			requestData <- word
		default:
		}
	}
}

type ByFavLink []*linkData

func (a ByFavLink) Len() int           { return len(a) }
func (a ByFavLink) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFavLink) Less(i, j int) bool { return a[i].Likes > a[j].Likes }

type ByRetLinks []*linkData

func (a ByRetLinks) Len() int           { return len(a) }
func (a ByRetLinks) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRetLinks) Less(i, j int) bool { return a[i].Retweets > a[j].Retweets }

func processLinks(lm linksMap) *linksDisplay {
	links := getLinksValues(lm)
	linksByFav := make(linksSlice, len(lm))
	linksByRet := make(linksSlice, len(lm))
	copy(linksByFav, links)
	copy(linksByRet, links)
	sort.Sort(ByFavLink(linksByFav))
	sort.Sort(ByFavLink(linksByRet))
	return &linksDisplay{linksByFav, linksByRet}
}

type ByFavTweet []*anaconda.Tweet

func (a ByFavTweet) Len() int           { return len(a) }
func (a ByFavTweet) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFavTweet) Less(i, j int) bool { return a[i].FavoriteCount > a[j].FavoriteCount }

type ByRetTweet []*anaconda.Tweet

func (a ByRetTweet) Len() int           { return len(a) }
func (a ByRetTweet) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRetTweet) Less(i, j int) bool { return a[i].RetweetCount > a[j].RetweetCount }

func processTweets(tweetsMap map[string]*anaconda.Tweet) *tweetsDisplay {
	// we will use map
	// type of max -> array of indexes
	// array of tweets
	// for some reason we need to use the ids because otherwise it appends to the end
	tweets := getTweetValues(tweetsMap)
	tweetsByFav := make([]*anaconda.Tweet, len(tweetsMap))
	tweetsByRet := make([]*anaconda.Tweet, len(tweetsMap))
	copy(tweetsByFav, tweets)
	copy(tweetsByRet, tweets)
	sort.Sort(ByFavTweet(tweetsByFav))
	sort.Sort(ByRetTweet(tweetsByRet))
	return &tweetsDisplay{tweetsByFav, tweetsByRet}
}

func arrayIndexes(len int) []int {
	result := make([]int, len)
	for i := 0; i < len; i++ {
		result[i] = i
	}
	return result
}
