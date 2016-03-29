package main

import (
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
	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/:word", func(c *gin.Context) {
		comm <- 1
		td := <-tweets
		// maybe move this logic in processor since there is not need for it here
		// or find out the way to sort by index
		// TODO fix sorting here
		tweetsByLikes := make([]*anaconda.Tweet, len(td.tweets))
		tweetsByRetweets := make([]*anaconda.Tweet, len(td.tweets))
		for i := range td.favInd {
			tweetsByLikes[i] = td.tweets[td.favInd[i]]
			tweetsByRetweets[i] = td.tweets[td.retwInd[i]]
		}
		word := c.Param("word")
		if isBeingTracked(word) {
			c.HTML(http.StatusOK, "word.tmpl", gin.H{
				"title":            "Main website",
				"tweetsByLikes":    tweetsByLikes[1:10],
				"tweetsByRetweets": tweetsByRetweets[1:10],
			})
		} else {
			c.String(http.StatusNotFound, "This words is not followed")
		}
	})
	router.Run()
}
