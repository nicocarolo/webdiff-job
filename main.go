package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gocron"
	"github.com/webdiff-job/src/db"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
		log.Println("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/ping", Ping)

	router.Run(":" + port)

	go func() {
		gocron.Every(1).Minute().Do(taskWithParams, 1, "hello")
		<-gocron.Start()
	}()
}

func Ping(c *gin.Context) {
	session, err := db.GetMongoSession()
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		c.JSON(599, gin.H{
			"message": "Can't connect to database",
		})
		return
	}
	defer db.CloseMongoSession(session)

	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func taskWithParams(a int, b string) {
	fmt.Println(a, b)
}
