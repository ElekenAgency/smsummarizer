package main

import (
	"sort"
)

func processor(reqMainC <-chan string, respMainC chan<- *displayData) {
	respDataC := make(chan *dataChannelValues)
	reqDataC := make(chan string)
	go dataManager(respDataC, reqDataC)
	for {
		select {
		case <-dReqC:
			return
		case tweetsAndLinks := <-respDataC:
			respMainC <- &displayData{tweets: processTweets(tweetsAndLinks.tweets),
				links: processLinks(tweetsAndLinks.links)}
		case word := <-reqMainC:
			reqDataC <- word
		default:
		}
	}
}

type byFavLink linksSlice

func (a byFavLink) Len() int           { return len(a) }
func (a byFavLink) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byFavLink) Less(i, j int) bool { return a[i].Likes > a[j].Likes }

type byRetLinks linksSlice

func (a byRetLinks) Len() int           { return len(a) }
func (a byRetLinks) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byRetLinks) Less(i, j int) bool { return a[i].Retweets > a[j].Retweets }

func processLinks(lm linksMap) *linksDisplay {
	links := getLinksValues(lm)
	linksByFav := make(linksSlice, len(lm))
	linksByRet := make(linksSlice, len(lm))
	copy(linksByFav, links)
	copy(linksByRet, links)
	sort.Sort(byFavLink(linksByFav))
	sort.Sort(byFavLink(linksByRet))
	return &linksDisplay{linksByFav, linksByRet}
}

type byFavTweet tweetsSlice

func (a byFavTweet) Len() int           { return len(a) }
func (a byFavTweet) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byFavTweet) Less(i, j int) bool { return a[i].FavoriteCount > a[j].FavoriteCount }

type byRetTweet tweetsSlice

func (a byRetTweet) Len() int           { return len(a) }
func (a byRetTweet) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byRetTweet) Less(i, j int) bool { return a[i].RetweetCount > a[j].RetweetCount }

func processTweets(tm tweetsMap) *tweetsDisplay {
	tweets := getTweetValues(tm)
	tweetsByFav := make(tweetsSlice, len(tm))
	tweetsByRet := make(tweetsSlice, len(tm))
	copy(tweetsByFav, tweets)
	copy(tweetsByRet, tweets)
	sort.Sort(byFavTweet(tweetsByFav))
	sort.Sort(byRetTweet(tweetsByRet))
	return &tweetsDisplay{tweetsByFav, tweetsByRet}
}

func arrayIndexes(len int) []int {
	result := make([]int, len)
	for i := 0; i < len; i++ {
		result[i] = i
	}
	return result
}
