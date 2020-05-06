package main

import (
	"encoding/json"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

//struct to store user credentials
type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

//struct for NASA images
type ImageSource struct {
	Source string `json:"img_src"`
}

//struct for response from NASA API
type Response struct {
	ImageList []ImageSource `json:"photos"`
}

func getClient(creds *Credentials) (*twitter.Client, error) {

	//credentials passed in from environment variables
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// verify credentials
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	user, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	log.Printf("User's ACCOUNT:\n%+v\n", user)
	return client, nil
}

func main() {

	fmt.Println("Rover Bot v1.0")

	//downloads random image from Mars and saves it to local dir as Mars.jpg
	downloadImage(randNASA())

	tweetImage()
	deleteImage()
}

func getImages(body []byte) (*Response, error) {
	var s = new(Response)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	return s, err
}

func randNASA() string {
	rand.Seed(time.Now().UnixNano())

	//gets and sets key for NASA API
	KEY := os.Getenv("NASA_KEY")
	response, err := http.Get("https://api.nasa.gov/mars-photos/api/v1/rovers/curiosity/photos?sol=1000&api_key=" + KEY)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	//gets images from NASA API and
	s, err := getImages(responseData)
	randNumber := rand.Intn(len(s.ImageList))
	randImage := s.ImageList[randNumber].Source
	return randImage
}

func tweetImage() {
	creds := Credentials{
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
	}

	client, err := getClient(&creds)
	if err != nil {
		log.Println("Error getting Twitter Client")
		log.Println(err)
	}

	//test tweet!
	tweet, resp, err := client.Statuses.Update(" I am putting myself to the fullest possible use, which is all I think that any conscious entity can ever hope to do.", nil)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%+v\n", resp)
	log.Printf("%+v\n", tweet)
}

func downloadImage(url string) {
	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	file, err := os.Create("Mars.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success!")
}

func deleteImage() {
	path := "Mars.jpg"
	err := os.Remove(path)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Mars.jpg deleted")
}
