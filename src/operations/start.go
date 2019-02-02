package operations

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/webdiff-job/src/db"
	"github.com/webdiff-job/src/job"
	"github.com/webdiff-job/src/models"
)

func Start(c *gin.Context) {
	session, err := db.GetMongoSession()
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		c.JSON(599, gin.H{
			"message": "Can't connect to database",
		})
		return
	}
	defer db.CloseMongoSession(session)

	collection := session.DB("heroku_rvdsxf5j").C("webs")

	var webs []models.Web

	err = collection.Find(nil).All(&webs)
	if err != nil {
		log.Println("Not exists Webs to inspect")
		return
	}

	notified := job.Inspect(webs, collection)
	log.Println(fmt.Sprintf("Finished start, notified %d webs", len(notified)))
}
