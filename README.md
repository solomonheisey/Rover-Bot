# Rover-Bot
Rover Bot is an Twitter bot I coded up primarily in Go. Rover Bot works by implementing 2 APIS. First, a random image is selected from the NASA Mars Rover API and downloaded, next a phrase is generated from a data-set gathered from Twitter. The Twitter data set contains over 10,000 cleansed phrases revolving under the subject matter of weather. After both the image and the caption are created, a connection to the Twitter API is made in Python and the JSON object is sent using REST API.

## Getting Started
```
$ pip install -r requirements.txt
$ go run main.go
```


 
 
