package form

import (
	"log"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/onihani/go-tweet/internal/utils"
)

func CollectStringInput(label string) (string, error) {
	var input string

	// Initialize the form
	form := huh.NewForm(
		// Add a field to collect the input
		huh.NewGroup(
			huh.NewInput().Title(label).Value(&input),
		),
	)

	// Run the form
	err := form.Run()
	if err != nil {
		log.Fatalf("Failed to get input: %v", err)
		return "", err
	}

	// Return the tweet URL
	return input, nil
}

func CollectTweetUrl() (string, error) {
	var tweetUrl string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter the tweet url").
				Value(&tweetUrl).
				Validate(utils.IsValidTwitterUrl),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatalf("Failed to get tweet url: %v", err)
		return "", err
	}

	return tweetUrl, nil
}

func GetDirectoryPath(label string, isDirectory bool) (string, error) {
	var input string

	if isDirectory {
		// Initialize file picker for directory selection
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Printf("Failed to get home directory: %v", err)
			homeDir = "."
		}

		filePicker := huh.NewFilePicker().
			Title(label).
			DirAllowed(true).
			FileAllowed(false).
			CurrentDirectory(homeDir).
			ShowHidden(true).
			// Height(30).
			Value(&input)

		// Initialize the form with file picker
		form := huh.NewForm(
			huh.NewGroup(filePicker),
		)

		// Run the form
		err = form.Run()
		if err != nil {
			log.Fatalf("Failed to get directory input: %v", err)
			return "", err
		}
	} else {
		// Original input field for non-directory inputs
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title(label).Value(&input),
			),
		)

		// Run the form
		err := form.Run()
		if err != nil {
			log.Fatalf("Failed to get input: %v", err)
			return "", err
		}
	}

	// Return the selected input or directory path
	return input, nil
}

func SelectResolution(resolutions []string) (string, error) {
	var options = make([]huh.Option[string], len(resolutions), cap(resolutions))

	for index, resolution := range resolutions {
		options[index] = huh.NewOption(resolution, resolution)
	}

	var selectedResolution string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose video resolution").
				Description("Select the resolution you would prefer to download this media at").
				Options(options...).
				Value(&selectedResolution),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatalf("Failed to get prefered resolution: %v", err)
		return "", err
	}

	return selectedResolution, nil
}
