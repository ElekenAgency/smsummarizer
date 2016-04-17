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

func GetMainEngine() *gin.Engine {
	tweets, comm := make(chan *TweetsData), make(chan interface{})
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
				tweetsByLikes, tweetsByRetweets := getStats(tweets, comm)
				if len(tweetsByLikes) > 10 {
					tweetsByLikes = tweetsByLikes[0:10]
					tweetsByRetweets = tweetsByRetweets[0:10]
				}
				c.HTML(http.StatusOK, "word.tmpl", gin.H{
					"title":            "Main website",
					"tweetsByLikes":    tweetsByLikes,
					"tweetsByRetweets": tweetsByRetweets,
				})
			} else {
				c.String(http.StatusNotFound, "This words is not followed")
			}
		})
	}
	api := r.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"title": "APIs",
				"words": trackingWords,
			})
		})
		api.GET("/:word", func(c *gin.Context) {
			word := c.Param("word")
			if isBeingTracked(word) {
				tweetsByLikes, tweetsByRetweets := getStats(tweets, comm)
				if len(tweetsByLikes) > 10 {
					tweetsByLikes = tweetsByLikes[0:10]
					tweetsByRetweets = tweetsByRetweets[0:10]
				}
				// TODO Probably add marshaling if needed
				tl := simplifyTweets(tweetsByLikes)
				tr := simplifyTweets(tweetsByRetweets)
				c.JSON(http.StatusOK, gin.H{
					"tweetsByLikes":    tl,
					"tweetsByRetweets": tr,
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
	Logger = initLog()
	GetMainEngine().Run(":5000")
}

func (tweet TweetShort) MarshalJSON() ([]byte, error) {
	return json.Marshal(tweet.text)
}
