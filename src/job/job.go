package job

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/webdiff-job/src/config"
	"github.com/webdiff-job/src/db"
	"github.com/webdiff-job/src/models"
)

func Job() {
	log.Println("Starting job")
	session, err := db.GetMongoSession()
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
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

	notified := Inspect(webs, collection)
	log.Println(fmt.Sprintf("Finished job, notified %d webs", len(notified)))
}

func Inspect(webs []models.Web, collection *mgo.Collection) map[string][]models.Web {
	log.Println("Start inspect")
	notified := map[string][]models.Web{}

	for _, web := range webs {
		client := &http.Client{}
		req, err := http.NewRequest("GET", web.Url, nil)
		req.Header.Add("user-agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36`)
		response, httperror := client.Do(req)
		//  := http.Get(web.Url)
		if httperror != nil {
			log.Println(fmt.Errorf("Error while requesting url: %s", httperror.Error()))
		}

		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println(fmt.Errorf("Error while reading url: %s", err.Error()))
		}
		response.Body.Close()

		checksum := sha256.Sum256(contents)
		strChecksum := hex.EncodeToString(checksum[:])

		if strChecksum != web.Web {
			// if web.Url == "https://www.jumbo.com.ar/almacen/desayuno-y-merienda/cafes?PS=18" {
			// 	// fmt.Println("NOTIFICAAAAAA")
			// 	// fmt.Println(string(contents[:]))
			// 	f, _ := os.Create("/tmp/web.txt")
			// 	f.Write(contents)
			// }
			log.Println(fmt.Sprintf("Updating %s from %s", web.Web, web.Url))
			web.Web = strChecksum
			err = collection.Update(bson.M{"_id": web.Id}, web)
			if err != nil {
				log.Println(fmt.Errorf("Error while update web: %s error: %s", web.Url, err.Error()))
			}
			webByIdNotified, wasNotified := notified[web.WebId]
			if !wasNotified {
				log.Println(fmt.Sprintf("Notifying %s from %s", web.Url, web.WebId))
				values := map[string]string{"merchant": web.WebId}
				jsonValue, _ := json.Marshal(values)
				response, httperror = http.Post(fmt.Sprintf(config.GetScrapperURL(), "/process"),
					"application/json", bytes.NewBuffer(jsonValue))
				if httperror != nil {
					log.Println(fmt.Errorf("Error while post to scrapper: %s", httperror.Error()))
				}
			} else {
				log.Println(fmt.Sprintf("Already notified %s from %s", web.Url, web.WebId))
			}
			webByIdNotified = append(webByIdNotified, web)
			notified[web.WebId] = webByIdNotified
		}
	}

	return notified
}
