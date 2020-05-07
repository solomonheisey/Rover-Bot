package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

//struct for NASA images
type ImageSource struct {
	Source string `json:"img_src"`
}

//struct for response from NASA API
type Response struct {
	ImageList []ImageSource `json:"photos"`
}

//struct for tweet content from json
type Tweets struct {
	Tweet string `json:"tweet_text"`
}

//struct for metadata from json
type MetaData struct {
	Data []Tweets `json:"tweets"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("Rover Bot v1.0")

	//downloads random image from Mars and saves it to local dir as Mars.jpg
	downloadImage(randNASA())
	tweetImage(randQuote())
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

//selects random quote from database
func randQuote() string {
	jsonFile, err := os.Open("weather.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var data MetaData

	json.Unmarshal(byteValue, &data)
	randNumber := rand.Intn(len(data.Data))

	removeTags := strings.Replace(data.Data[randNumber].Tweet, "{link}", "", -1)
	removeTags = strings.Replace(removeTags, "@mention ", "", -1)
	removeTags = strings.Replace(removeTags, "RT @mention: ", "", -1 )

	return removeTags
}

//Tweets image with caption...also pipes Python output to Go
func tweetImage(message string) {
	cmd := exec.Command("python", "image.py", message)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	go copyOutput(stdout)
	go copyOutput(stderr)
	cmd.Wait()
}

//copies output from Python script
func copyOutput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

//downloads random image from Mars rover
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
}

//deleted downloaded file
func deleteImage() {
	path := "Mars.jpg"
	err := os.Remove(path)

	if err != nil {
		fmt.Println(err)
		return
	}
}
