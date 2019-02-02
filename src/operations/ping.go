package operations

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/webdiff-job/src/db"
)

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
