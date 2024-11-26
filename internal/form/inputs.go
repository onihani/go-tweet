package form

import (
	"log"

	"github.com/charmbracelet/huh"
)

// CollectTweetURL asks the user to input a tweet URL via a form
func CollectTweetURL() (string, error) {
	var tweetUrl string

	// Initialize the form
	form := huh.NewForm(
		// Add a field to collect the tweet URL
		huh.NewGroup(
			huh.NewInput().Title("Enter the tweet URL").Value(&tweetUrl),
		),
	)

	// Run the form
	err := form.Run()
	if err != nil {
		log.Fatalf("Failed to get input: %v", err)
		return "", err
	}

	// Return the tweet URL
	return tweetUrl, nil
}
