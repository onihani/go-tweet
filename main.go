package main

import (
	"fmt"

	"github.com/charmbracelet/huh/spinner"
	"github.com/onihani/go-tweet/internal/downloader"
	"github.com/onihani/go-tweet/internal/form"
	"github.com/onihani/go-tweet/internal/models"
	"github.com/onihani/go-tweet/internal/utils"
)

func main() {
	url, err := form.CollectTweetUrl()
	if err != nil {
		fmt.Println("Failed to get tweet url:", err)
		return
	}

	var tweetId string
	err = spinner.New().Title("Extracting tweet id from url...").Accessible(false).Action(func() {
		id, err := utils.ExtractTweetID(url)
		if err != nil {
			panic(err)
		}

		tweetId = id
	}).Run()

	if err != nil {
		fmt.Println("Failed to extract tweet id from url:", err)
		return
	}

	var tweet *models.Tweet
	err = spinner.New().Title("Fetching tweet data...").Accessible(false).Action(func() {
		t, err := downloader.FetchTweet(tweetId)
		if err != nil {
			panic(err)
		}

		tweet = t
	}).Run()

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

	selectedResolution, err := form.SelectResolution(resolutions)
	if err != nil {
		fmt.Println("Error selecting resolutin:", err)
		return
	}

	selectedVideo := tweetVideos[selectedResolution]

	downloadDirectory, err := form.CollectStringInput("Enter path to the directory you want to download video to")
	if err != nil {
		fmt.Println("Error collecting download directory:", err)
		return
	}

	var downloadedFilePath string
	err = spinner.New().
		Title(fmt.Sprintf("Downloading video at %s to %s...", selectedResolution, downloadDirectory)).
		Accessible(false).
		Action(func() {
			dfPath, err := selectedVideo.Download(downloadDirectory, tweet.Text)
			if err != nil {
				panic(err)
			}

			downloadedFilePath = dfPath
		}).Run()

	if err != nil {
		fmt.Println("Failed to download video", err)
		return
	}

	fmt.Println("Video downloaded to:", downloadedFilePath)
}
