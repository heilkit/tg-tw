package twitter

import (
	"fmt"
	"regexp"
	"strings"
)

var VxRegex = regexp.MustCompile("https://api\\.vxtwitter\\.com/[a-zA-Z0-9_]{1,15}/status/[0-9]+")

func Vx(url string) (string, error) {
	if !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	switch {
	case strings.HasPrefix(url, "https://twitter.com"):
		url = strings.Replace(url, "https://twitter.com", "https://api.vxtwitter.com", 1)
	case strings.HasPrefix(url, "https://x.com"):
		url = strings.Replace(url, "https://x.com", "https://api.vxtwitter.com", 1)
	case strings.HasPrefix(url, "https://vxtwitter.com"):
		url = strings.Replace(url, "https://vxtwitter.com", "https://api.vxtwitter.com", 1)
	default:
		return "", fmt.Errorf("url has to be 'twitter.com/*' (got '%s')", url)
	}

	parsed := VxRegex.FindAll([]byte(url), 1)
	if len(parsed) == 0 {
		return "", fmt.Errorf("url parsing failed")
	}
	return string(parsed[0]), nil
}

type VxMedia struct {
	AltText        string `json:"altText"`
	DurationMillis int    `json:"duration_millis,omitempty"`
	Size           struct {
		Height int `json:"height"`
		Width  int `json:"width"`
	} `json:"size"`
	ThumbnailUrl string `json:"thumbnail_url"`
	Type         string `json:"type"`
	Url          string `json:"url"`
}

type VxPost struct {
	Date           string    `json:"date"`
	DateEpoch      int       `json:"date_epoch"`
	Hashtags       []string  `json:"hashtags"`
	Likes          int       `json:"likes"`
	MediaURLs      []string  `json:"mediaURLs"`
	MediaExtended  []VxMedia `json:"media_extended"`
	Replies        int       `json:"replies"`
	Retweets       int       `json:"retweets"`
	Text           string    `json:"text"`
	TweetID        string    `json:"tweetID"`
	TweetURL       string    `json:"tweetURL"`
	UserName       string    `json:"user_name"`
	UserScreenName string    `json:"user_screen_name"`
}
