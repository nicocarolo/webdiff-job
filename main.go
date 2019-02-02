package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocron"
	"github.com/webdiff-job/src/job"
	"github.com/webdiff-job/src/operations"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "4000"
		log.Println("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/ping", operations.Ping)

	router.POST("/add", operations.Add)

	router.POST("/start", operations.Start)

	go func() {
		gocron.Every(3).Minutes().Do(job.Job)
		<-gocron.Start()
	}()

	router.Run(":" + port)
}

func taskWithParams(a int, b string) {
	fmt.Println(time.Now().Format(time.RFC850))
	fmt.Println(a, b)
}
