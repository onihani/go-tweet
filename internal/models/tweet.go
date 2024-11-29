package models

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/onihani/go-tweet/internal/utils"
)

type Tweet struct {
	TypeName          string         `json:"__typename"`
	Lang              string         `json:"lang"`
	FavoriteCount     int            `json:"favorite_count"`
	PossiblySensitive bool           `json:"possibly_sensitive"`
	CreatedAt         string         `json:"created_at"`
	DisplayTextRange  []int          `json:"display_text_range"`
	Entities          TweetEntities  `json:"entities"`
	ID                string         `json:"id_str"`
	Text              string         `json:"text"`
	User              TweetUser      `json:"user"`
	EditControl       EditControl    `json:"edit_control"`
	MediaDetails      []MediaDetails `json:"mediaDetails"`
	Photos            []interface{}  `json:"photos"`
	Video             VideoDetails   `json:"video"`
	ConversationCount int            `json:"conversation_count"`
	NewsActionType    string         `json:"news_action_type"`
	IsEdited          bool           `json:"isEdited"`
	IsStaleEdit       bool           `json:"isStaleEdit"`
}

type TweetEntities struct {
	Hashtags     []interface{} `json:"hashtags"`
	Urls         []interface{} `json:"urls"`
	UserMentions []interface{} `json:"user_mentions"`
	Symbols      []interface{} `json:"symbols"`
	Media        []TweetMedia  `json:"media"`
}

type TweetMedia struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type TweetUser struct {
	ID                string `json:"id_str"`
	Name              string `json:"name"`
	ProfileImageURL   string `json:"profile_image_url_https"`
	ScreenName        string `json:"screen_name"`
	Verified          bool   `json:"verified"`
	IsBlueVerified    bool   `json:"is_blue_verified"`
	ProfileImageShape string `json:"profile_image_shape"`
}

type EditControl struct {
	EditTweetIDs       []string `json:"edit_tweet_ids"`
	EditableUntilMsecs string   `json:"editable_until_msecs"`
	IsEditEligible     bool     `json:"is_edit_eligible"`
	EditsRemaining     string   `json:"edits_remaining"`
}

type MediaDetails struct {
	AdditionalMediaInfo map[string]interface{} `json:"additional_media_info"`
	DisplayURL          string                 `json:"display_url"`
	ExpandedURL         string                 `json:"expanded_url"`
	MediaURL            string                 `json:"media_url_https"`
	Type                string                 `json:"type"`
	URL                 string                 `json:"url"`
	VideoInfo           VideoInfo              `json:"video_info"`
	OriginalInfo        OriginalInfo           `json:"original_info"`
}

type VideoInfo struct {
	AspectRatio []int          `json:"aspect_ratio"`
	DurationMs  int            `json:"duration_millis"`
	Variants    []VideoVariant `json:"variants"`
}

type VideoVariant struct {
	ContentType string `json:"content_type"`
	URL         string `json:"url"`
	Bitrate     int    `json:"bitrate,omitempty"`
}

type OriginalInfo struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

type VideoDetails struct {
	AspectRatio       []int  `json:"aspectRatio"`
	ContentType       string `json:"contentType"`
	DurationMs        int    `json:"durationMs"`
	MediaAvailability struct {
		Status string `json:"status"`
	} `json:"mediaAvailability"`
	Poster   string         `json:"poster"`
	Variants []VideoVariant `json:"variants"`
	VideoId  struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"videoId"`
	ViewCount int `json:"viewCount"`
}

func (t *Tweet) GetImages() []string {
	var imageURLs []string

	for _, media := range t.MediaDetails {
		if media.Type == "photo" {
			imageURLs = append(imageURLs, media.MediaURL)
		}
	}

	return imageURLs
}

func (t *Tweet) GetVideos() map[string]VideoVariant {
	videoMap := make(map[string]VideoVariant)

	// Regular expression to extract resolution from the video URL (e.g., "320x562").
	resolutionRegex := regexp.MustCompile(`(\d+x\d+)`)

	// Iterate through the media details
	for _, media := range t.MediaDetails {
		if media.Type == "video" || media.Type == "animated_gif" {
			// Iterate through video variants
			for _, variant := range media.VideoInfo.Variants {
				// Only consider "video/mp4" type variants
				if variant.ContentType != "video/mp4" {
					continue
				}

				// Extract resolution from the URL
				matches := resolutionRegex.FindStringSubmatch(variant.URL)
				if len(matches) < 2 {
					// If resolution is not found, skip this variant
					continue
				}

				resolution := matches[1]
				videoMap[resolution] = variant
			}
		}
	}

	return videoMap
}

func (t *Tweet) GetAllMedia() []string {
	var mediaURLs []string

	for _, media := range t.MediaDetails {
		if media.Type == "photo" {
			mediaURLs = append(mediaURLs, media.MediaURL)
		} else if media.Type == "video" || media.Type == "animated_gif" {
			for _, videoVariant := range media.VideoInfo.Variants {
				mediaURLs = append(mediaURLs, videoVariant.URL)
			}
		}
	}

	return mediaURLs
}

func (v *VideoVariant) Download(destinationDir string, title string) (string, error) {
	// Expand `~` to the user's home directory
	if strings.HasPrefix(destinationDir, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		destinationDir = filepath.Join(homeDir, destinationDir[2:])
	}

	// Ensure the directory exists or create it
	if err := os.MkdirAll(destinationDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Create a file name based on bitrate and aspect ratio
	fileName := fmt.Sprintf("%s_%dkbp_%d.mp4", utils.SanitizeTitle(title), v.Bitrate/1000, time.Now().UnixMilli())
	filePath := filepath.Join(destinationDir, fileName)

	// Open the URL
	resp, err := http.Get(v.URL)
	if err != nil {
		return "", fmt.Errorf("failed to download video: %w", err)
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	// Create the destination file
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write the response body to the file
	if _, err := io.Copy(file, resp.Body); err != nil {
		return "", fmt.Errorf("failed to save video: %w", err)
	}

	return filePath, nil
}
