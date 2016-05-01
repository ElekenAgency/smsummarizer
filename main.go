package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
)

var respC chan *displayData
var reqC chan string

func cleanup() {
	fmt.Println("\nExiting!")
	// send 2 request to tell both processor and data to end
	dReqC <- 1
	dReqC <- 1
	data := <-dRespC
	writeDumpContents(tweetsDumpPath, data.tweets)
	writeDumpContents(linksDumpPath, data.links)
	close(dReqC)
	close(dRespC)
	close(respC)
	close(reqC)
}

func init() {
	flag.Var(&trackedWords, "words", "Words to track")
	// setup listening to CTRL-C and SIGTERM that docker send
	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGINT,
		syscall.SIGTERM)
	summarizerDumpPath := os.Getenv("SM_DUMP")
	if summarizerDumpPath == "" {
		summarizerDumpPath = "."
	}
	if _, err := os.Stat(summarizerDumpPath); err == nil {
		tweetsDumpPath = path.Join(summarizerDumpPath, "dump_tweets")
		linksDumpPath = path.Join(summarizerDumpPath, "dump_links")
	} else {
		fmt.Println("Specified dump path doesn't exist")
	}
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()
}

func isBeingTracked(word string) bool {
	for _, trackedWord := range trackedWords {
		if word == trackedWord {
			return true
		}
	}
	return false
}

func getDispayData(word string, respC chan *displayData, reqC chan string) *displayData {
	reqC <- word
	data := <-respC
	if len(data.tweets.tweetsByFav) > 10 {
		data.tweets.tweetsByFav = data.tweets.tweetsByFav[0:10]
		data.tweets.tweetsByRet = data.tweets.tweetsByRet[0:10]
	}
	if len(data.links.linksByFav) > 10 {
		data.links.linksByFav = data.links.linksByFav[0:10]
		data.links.linksByRet = data.links.linksByRet[0:10]
	}
	return data
}

// GetMainEngine creates a git.Engine with routes
// Exported for testing
func GetMainEngine() *gin.Engine {
	respC, reqC = make(chan *displayData), make(chan string)
	dReqC, dRespC = make(dReqChan), make(dRespChan)
	go processor(reqC, respC)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	// index router
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
			"words": trackedWords,
		})
	})
	web := r.Group("/web")
	{
		web.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"title": "Web",
				"words": trackedWords,
			})
		})
		web.GET("/:word", func(c *gin.Context) {
			word := c.Param("word")
			if isBeingTracked(word) {
				data := getDispayData(word, respC, reqC)
				c.HTML(http.StatusOK, "word.tmpl", gin.H{
					"title":            "Main website",
					"tweetsByLikes":    data.tweets.tweetsByFav,
					"tweetsByRetweets": data.tweets.tweetsByRet,
					"linksByLikes":     data.links.linksByFav,
					"linksByRetweets":  data.links.linksByRet,
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
