package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/gocron"
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

	gocron.Every(1).Minute().Do(taskWithParams, 1, "hello")
}

func Ping(c *gin.Context) {
	env := os.Getenv("ENVIRONMENT")
	var url string

	if env == "PRODUCTION" {
		url = "mongodb://api:dM6CYayNQu8qr9b@ds149984.mlab.com:49984/heroku_rjnls62m"
	} else {
		url = "localhost"
	}
	session, err := mgo.Dial(url)

	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		c.JSON(599, gin.H{
			"message": "Can't connect to database",
		})
		return
	}
	session.SetSafe(&mgo.Safe{})
	defer session.Close()

	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func taskWithParams(a int, b string) {
	fmt.Println(a, b)
}
