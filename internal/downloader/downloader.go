package downloader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/onihani/go-tweet/internal/constants"
	"github.com/onihani/go-tweet/internal/models"
	"github.com/onihani/go-tweet/internal/utils"
)

func FetchTweet(id string) (*models.Tweet, error) {
	baseURL := fmt.Sprintf("%s/tweet-result", constants.TWITTER_SYNDICATION_URL)
	token := utils.GetToken(id)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s?id=%s&token=%s", baseURL, id, token), nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch tweet")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var tweet models.Tweet
	if err := json.Unmarshal(body, &tweet); err != nil {
		return nil, err
	}

	return &tweet, nil
}
