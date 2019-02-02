package operations

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/webdiff-job/src/config"
	"github.com/webdiff-job/src/db"
	"github.com/webdiff-job/src/models"
)

func Add(c *gin.Context) {
	session, err := db.GetMongoSession()
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		c.JSON(599, gin.H{
			"message": "Can't connect to database",
		})
		return
	}
	defer db.CloseMongoSession(session)

	collection := session.DB(config.WebDB).C(config.WebCollection)

	var request models.Request
	err = c.BindJSON(&request)

	var result models.Web
	err = collection.Find(bson.M{"url": request.Url}).One(&result)

	if err == nil {
		log.Println(fmt.Errorf("The requested url %s already exists", request.Url))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "The requested url already exists",
		})
		return
	}

	response, httperror := http.Get(request.Url)
	if httperror != nil {
		log.Println(fmt.Errorf("Error while requesting url: %s", err.Error()))
		c.JSON(http.StatusGatewayTimeout, gin.H{
			"error": "Error while requesting url",
		})
		return
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(fmt.Errorf("Error while reading url: %s", err.Error()))
		c.JSON(http.StatusGatewayTimeout, gin.H{
			"error": "Error while reading url",
		})
		return
	}

	checksum := sha256.Sum256(contents)
	err = collection.Insert(&models.Web{
		WebId:          request.Id,
		Url:            request.Url,
		Web:            hex.EncodeToString(checksum[:]),
		LastDateUpdate: strconv.FormatInt(time.Now().Unix(), 10),
	})
	if err != nil {
		log.Println(fmt.Errorf("Error while insert url: %s", err.Error()))
		c.JSON(http.StatusGatewayTimeout, gin.H{
			"error": "Error while insert url",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Added url: " + request.Url,
		"web":     checksum,
	})
}
