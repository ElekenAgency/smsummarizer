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

type ByFavLink linksSlice

func (a ByFavLink) Len() int           { return len(a) }
func (a ByFavLink) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFavLink) Less(i, j int) bool { return a[i].Likes > a[j].Likes }

type ByRetLinks linksSlice

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

type ByFavTweet tweetsSlice

func (a ByFavTweet) Len() int           { return len(a) }
func (a ByFavTweet) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFavTweet) Less(i, j int) bool { return a[i].FavoriteCount > a[j].FavoriteCount }

type ByRetTweet tweetsSlice

func (a ByRetTweet) Len() int           { return len(a) }
func (a ByRetTweet) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRetTweet) Less(i, j int) bool { return a[i].RetweetCount > a[j].RetweetCount }

func processTweets(tm tweetsMap) *tweetsDisplay {
	tweets := getTweetValues(tm)
	tweetsByFav := make(tweetsSlice, len(tm))
	tweetsByRet := make(tweetsSlice, len(tm))
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
