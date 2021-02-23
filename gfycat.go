package gfycat

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// ClientConfig describes configuration options
type ClientConfig struct {
	ClientID string
	ClientSecret string
}

// Client describes the gfycat client
type Client struct {
	AuthResponse
	ClientCredentials
}

// Gfycat describes a Gfycat api response.
type Response struct {
	GfyItem struct {
		GfyID              string        `json:"gfyId"`
		GfyName            string        `json:"gfyName"`
		GfyNumber          string        `json:"gfyNumber"`
		WebmURL            string        `json:"webmUrl"`
		GifURL             string        `json:"gifUrl"`
		MobileURL          string        `json:"mobileUrl"`
		MobilePosterURL    string        `json:"mobilePosterUrl"`
		MiniURL            string        `json:"miniUrl"`
		MiniPosterURL      string        `json:"miniPosterUrl"`
		PosterURL          string        `json:"posterUrl"`
		Thumb100PosterURL  string        `json:"thumb100PosterUrl"`
		Max5MbGif          string        `json:"max5mbGif"`
		Max2MbGif          string        `json:"max2mbGif"`
		Max1MbGif          string        `json:"max1mbGif"`
		Gif100Px           string        `json:"gif100px"`
		Width              int           `json:"width"`
		Height             int           `json:"height"`
		AvgColor           string        `json:"avgColor"`
		FrameRate          int           `json:"frameRate"`
		NumFrames          int           `json:"numFrames"`
		Mp4Size            int           `json:"mp4Size"`
		WebmSize           int           `json:"webmSize"`
		GifSize            int           `json:"gifSize"`
		Source             int           `json:"source"`
		CreateDate         int           `json:"createDate"`
		Nsfw               string        `json:"nsfw"`
		Mp4URL             string        `json:"mp4Url"`
		Likes              string        `json:"likes"`
		Published          int           `json:"published"`
		Dislikes           string        `json:"dislikes"`
		ExtraLemmas        string        `json:"extraLemmas"`
		Md5                string        `json:"md5"`
		Views              int           `json:"views"`
		Tags               []string      `json:"tags"`
		UserName           string        `json:"userName"`
		Title              string        `json:"title"`
		Description        string        `json:"description"`
		LanguageText       string        `json:"languageText"`
		LanguageCategories string         `json:"languageCategories"`
		Subreddit          string        `json:"subreddit"`
		RedditID           string        `json:"redditId"`
		RedditIDText       string        `json:"redditIdText"`
		DomainWhitelist    []interface{} `json:"domainWhitelist"`
	} `json:"gfyItem"`
}

// New creates a Gfycat client.
func New(config ClientConfig) (Client, error) {
	return Client{
		ClientCredentials: ClientCredentials{
			config.ClientID,
			config.ClientSecret,
		},
	}, nil
}

// Authenticate will authenticate using client credentials. Returns a bearer token.
func (client *Client) Authenticate() error {
	authClient := NewAuthClient(client.ClientID, client.ClientSecret)
	// Defaults to getting client credentials
	authResponse, err := authClient.GetAccessToken()
	if err != nil {
		return err
	}

	client.AuthResponse = authResponse
	return nil
}

// CheckToken checks if AccessToken is set on the Gfycat client
func (client *Client) CheckToken() bool {
	if client.AccessToken == "" {
		return false
	} else {
		return true
	}
}

// GetGfycat gets a single gfycat by ID
func (client Client) GetGfycat(gfyID string) (Response, error) {
	if ok := client.CheckToken(); !ok {
		return Response{}, errors.New("Access token missing")
	}
	bearer := "Bearer "+client.AccessToken

	req, err := http.NewRequest("GET", GFYCATS_URL+"/"+gfyID, nil)
	if err != nil {
		return Response{}, err
	}
	req.Header.Add("Authorization", bearer)

	httpClient := &http.Client{}

	var resp *http.Response
	resp, err = httpClient.Do(req)
	if err != nil {
		return Response{}, err
	}
	if resp.StatusCode == 404 {
		var res map[string]interface{}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return Response{}, err
		}
		if err := json.Unmarshal(body, &res); err != nil {
			return Response{}, err
		}
		return Response{}, errors.New(res["errorMessage"].(string))
	}
	if resp.StatusCode == 403 {
		return Response{}, errors.New("unauthorized")
	}

	var gfycatResponse Response
	err = json.NewDecoder(resp.Body).Decode(&gfycatResponse)
	if err != nil {
		return Response{}, err
	}

	return gfycatResponse, nil
}