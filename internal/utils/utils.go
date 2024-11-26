package utils

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func GetToken(id string) string {
	num, _ := strconv.ParseFloat(id, 64)
	token := (num / 1e15) * math.Pi
	tokenStr := strings.ReplaceAll(strconv.FormatFloat(token, 'f', -1, 64), "0", "")
	return strings.ReplaceAll(tokenStr, ".", "")
}

// ExtractTweetID takes a tweet URL and returns the tweet ID
func ExtractTweetID(url string) (string, error) {
	// Remove any trailing spaces from the URL
	url = strings.TrimSpace(url)

	// Regular expression to match Twitter tweet URL structure
	re := regexp.MustCompile(`https://x\.com/[^/]+/status/(\d+)`)

	// Match the URL against the regular expression
	matches := re.FindStringSubmatch(url)

	if len(matches) > 1 {
		// Return the tweet ID
		return matches[1], nil
	}

	// If no match found, return an error
	return "", fmt.Errorf("invalid tweet URL")
}
