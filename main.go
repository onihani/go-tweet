package main

import (
	"encoding/json"
	"fmt"

	"github.com/onihani/go-tweet/internal/downloader"
	"github.com/onihani/go-tweet/internal/form"
	"github.com/onihani/go-tweet/internal/utils"
)

func main() {
	url, err := form.CollectTweetURL()
	if err != nil {
		fmt.Println("Failed to get tweet url:", err)
		return
	}

	tweetId, err := utils.ExtractTweetID(url)
	if err != nil {
		fmt.Println("Failed to extract tweet id from url:", err)
		return
	}

	tweet, err := downloader.FetchTweet(tweetId)
	if err != nil {
		fmt.Println("Error fetching tweet:", err)
		return
	}

	tweetVideos := tweet.GetVideos()

	var resolutions []string
	for resolution := range tweetVideos {
		resolutions = append(resolutions, resolution)
	}

	if len(resolutions) <= 0 {
		fmt.Println("No videos found")
		return
	}

	firstVid := tweetVideos[resolutions[1]]

	jsonBytesArray, err := json.MarshalIndent(firstVid, "", "  ")
	if err != nil {
		fmt.Println("Error preparing tweet json :", err)
		return
	}
	fmt.Printf("Tweet: %+v\n", string(jsonBytesArray))

	downloadedFilePath, err := firstVid.Download("downloads")
	if err != nil {
		fmt.Println("Failed to download", err)
		return
	}

	fmt.Println("Video downloaded to:", downloadedFilePath)
}
