package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type dumpRequestChan chan interface{}
type dumpResponceChan chan *dump

var dumpReq dumpRequestChan
var dumpRes dumpResponceChan

func cleanup() {
	fmt.Println("\nExiting!")
	dumpReq <- 1
	tweetsAndLinks := <-dumpRes
	jsonVal, err := json.Marshal(tweetsAndLinks.tweets)
	err = ioutil.WriteFile("/tweets/dump_tweets", jsonVal, 0644)
	if err != nil {
		fmt.Println("Problems with saving the data")
	}

	jsonVal, err = json.Marshal(tweetsAndLinks.links)
	err = ioutil.WriteFile("/tweets/dump_links", jsonVal, 0644)
	if err != nil {
		fmt.Println("Problems with saving the data")
	}
}

func init() {
	flag.Var(&trackingWords, "words", "Words to track")
	// setup listening to CTRL-C and SIGTERM that docker send
	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGINT,
		syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()
}

func isBeingTracked(word string) bool {
	for _, trackedWord := range trackingWords {
		if word == trackedWord {
			return true
		}
	}
	return false
}

type TweetShort struct {
	text     string
	likes    int
	retweets int
}

func getStats(word string, process chan *displayData, comm chan string) *displayData {
	comm <- word
	d := <-process
	if len(d.tweets.tweetsByFav) > 10 {
		d.tweets.tweetsByFav = d.tweets.tweetsByFav[0:10]
		d.tweets.tweetsByRet = d.tweets.tweetsByRet[0:10]
	}
	if len(d.links.linksByFav) > 10 {
		d.links.linksByFav = d.links.linksByFav[0:10]
		d.links.linksByRet = d.links.linksByRet[0:10]
	}
	return d
}

func GetMainEngine() *gin.Engine {
	tweets, comm := make(chan *displayData), make(chan string)
	dumpReq, dumpRes = make(dumpRequestChan), make(dumpResponceChan)
	go processor(comm, tweets)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	// index router
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
			"words": trackingWords,
		})
	})
	web := r.Group("/web")
	{
		web.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"title": "Web",
				"words": trackingWords,
			})
		})
		web.GET("/:word", func(c *gin.Context) {
			word := c.Param("word")
			if isBeingTracked(word) {
				d := getStats(word, tweets, comm)
				c.HTML(http.StatusOK, "word.tmpl", gin.H{
					"title":            "Main website",
					"tweetsByLikes":    d.tweets.tweetsByFav,
					"tweetsByRetweets": d.tweets.tweetsByRet,
					"linksByLikes":     d.links.linksByFav,
					"linksByRetweets":  d.links.linksByRet,
				})
			} else {
				c.String(http.StatusNotFound, "This words is not followed")
			}
		})
	}
	return r
}

func main() {
	flag.Parse()
	GetMainEngine().Run(":5000")
}
