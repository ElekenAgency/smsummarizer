package main

import (
	"flag"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/gin-gonic/gin"
	"log"
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

func main() {
	flag.Parse()
	Logger = initLog()
	tweets, comm := make(chan *anaconda.Tweet), make(chan int)
	go dataManager(trackingWords, Logger, tweetsNumber, tweets, comm)
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		comm <- 1
		tweet := <-tweets
		c.JSON(200, gin.H{
			"message": tweet.Text,
		})
	})
	router.Run()
}
