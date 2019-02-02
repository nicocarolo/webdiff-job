package config

import "os"

const WebDB = "heroku_rvdsxf5j"
const WebCollection = "webs"

func GetURLDB() string {
	env := os.Getenv("ENVIRONMENT")
	var url string

	if env == "PRODUCTION" {
		url = "mongodb://api:dM6CYayNQu8qr9b@ds147003.mlab.com:47003/heroku_rvdsxf5j"
	} else {
		url = "localhost"
	}
	return url
}

func GetScrapperURL() string {
	env := os.Getenv("ENVIRONMENT")
	var url string

	if env == "PRODUCTION" {
		url = "https://go-price-scrapper.herokuapp.com/%s"
	} else {
		url = "localhost:3000/%s"
	}
	return url
}
