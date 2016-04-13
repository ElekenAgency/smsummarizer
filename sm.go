package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var Logger *log.Logger

func initLog() *log.Logger {
	if *debugingMode {
		return log.New(os.Stdout, "DEBUG:", log.Ldate|log.Ltime|log.Lshortfile)
	}
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}

	return log.New(file,
		"PREFIX: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func cleanup() {
	fmt.Println("\nExiting!")
}

func init() {
	flag.Var(&trackingWords, "words", "Words to track")
	// setup listening to CTRL-C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
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

func getStats(tweets chan *TweetsData, comm chan interface{}) ([]*anaconda.Tweet, []*anaconda.Tweet) {
	comm <- 1
	td := <-tweets
	return td.tweetsByFav, td.tweetsByRet
}

func main() {
	flag.Parse()
	Logger = initLog()
	tweets, comm := make(chan *TweetsData), make(chan interface{})
	go processor(comm, tweets)

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	// index router
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
			"words": trackingWords,
		})
	})
	router.GET("/api", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "APIs",
			"words": trackingWords,
		})
	})
	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/web/:word", func(c *gin.Context) {
		word := c.Param("word")
		if isBeingTracked(word) {
			tweetsByLikes, tweetsByRetweets := getStats(tweets, comm)
			c.HTML(http.StatusOK, "word.tmpl", gin.H{
				"title":            "Main website",
				"tweetsByLikes":    tweetsByLikes[0:9],
				"tweetsByRetweets": tweetsByRetweets[0:9],
			})
		} else {
			c.String(http.StatusNotFound, "This words is not followed")
		}
	})
	router.GET("/api/:word", func(c *gin.Context) {
		word := c.Param("word")
		if isBeingTracked(word) {
			tweetsByLikes, tweetsByRetweets := getStats(tweets, comm)
			tl := simplifyTweets(tweetsByLikes[1:10])
			fmt.Print(tl)
			tr := simplifyTweets(tweetsByRetweets[1:10])
			tlj, err := json.Marshal(tl)
			fmt.Print(tlj)
			if err != nil {
				fmt.Println("error:", err)
			}
			trj, err := json.Marshal(tr)
			if err != nil {
				fmt.Println("error:", err)
			}
			c.JSON(http.StatusOK, gin.H{
				"tweetsByLikes":    trj,
				"tweetsByRetweets": tlj,
			})
		} else {
			c.String(http.StatusNotFound, "This words is not followed")
		}
	})
	router.Run()
}
