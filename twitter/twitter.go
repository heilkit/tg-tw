package twitter

import (
	"encoding/json"
	"fmt"
	"github.com/cavaliergopher/grab/v3"
	"net/http"
	"os"
	"sync"
	"time"
)

type API struct {
	Timeout time.Duration
	Sync    *sync.Mutex
	Client  *http.Client
}

func New() *API {
	return &API{
		Timeout: time.Second * 5,
		Sync:    &sync.Mutex{},
		Client:  http.DefaultClient,
	}
}

func (api *API) Get(url string) (*VxPost, error) {
	if api.Sync != nil {
		api.Sync.Lock()
		defer time.AfterFunc(api.Timeout, api.Sync.Unlock)
	}

	resp, err := api.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http.get '%s': %v", url, err)
	}
	defer resp.Body.Close()

	var post VxPost
	if err := json.NewDecoder(resp.Body).Decode(&post); err != nil {
		return nil, fmt.Errorf("json.decode '%s': %v", url, err)
	}

	return &post, nil
}

func (api *API) DownloadTempVx(url string) (files []string, dir string, post *VxPost, err error) {
	post, err = api.Get(url)
	if err != nil {
		return nil, "", nil, fmt.Errorf("api.download.get: %v", err)
	}

	dir, err = os.MkdirTemp("", "twitter_media_*")
	if err != nil {
		return nil, "", post, fmt.Errorf("api.download.mkdir: %v", err)
	}

	ret := []string{}
	client := grab.NewClient()
	client.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.3.1 Safari/605.1.1"
	client.HTTPClient = api.Client
	for _, mediaURL := range post.MediaURLs {
		mediaURL += "?name=large"
		req, err := grab.NewRequest(dir, mediaURL)
		if err != nil {
			return ret, dir, post, fmt.Errorf("api.download.grab.newrequest for '%s': %v", mediaURL, err)
		}

		resp := client.Do(req)
		if err := resp.Err(); err != nil {
			return ret, dir, post, fmt.Errorf("api.download.grab.dorequest for '%s': %v", mediaURL, err)
		}
		ret = append(ret, resp.Filename)
	}

	return ret, dir, post, nil
}

func (api *API) DownloadTemp(url string) (files []string, dir string, post *VxPost, err error) {
	parsed, err := Vx(url)
	if err != nil {
		return nil, "", nil, fmt.Errorf("api.download.parseurl for '%s': %v", url, err)
	}
	return api.DownloadTempVx(parsed)
}
