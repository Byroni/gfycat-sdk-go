package gfycat

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// ClientConfig describes configuration options
type ClientConfig struct {
	ClientID     string
	ClientSecret string
}

// GfycatClient describes the gfycat client
type GfycatClient struct {
	AuthResponse
	ClientCredentials
}

// Gfycat describes a Gfycat api response.
type GfycatResponse struct {
	GfyItem struct {
		GfyItem
	} `json:"gfyItem"`
}

// GfyItem describes the base gfycat item
type GfyItem struct {
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
	Width              json.Number   `json:"width"`
	Height             json.Number   `json:"height"`
	AvgColor           string        `json:"avgColor"`
	FrameRate          float32       `json:"frameRate"`
	NumFrames          float32       `json:"numFrames"`
	Mp4Size            json.Number   `json:"mp4Size"`
	WebmSize           json.Number   `json:"webmSize"`
	GifSize            json.Number   `json:"gifSize"`
	Source             json.Number   `json:"source"`
	CreateDate         json.Number   `json:"createDate"`
	Nsfw               json.Number   `json:"nsfw"`
	Mp4URL             string        `json:"mp4Url"`
	Likes              json.Number   `json:"likes"`
	Published          json.Number   `json:"published"`
	Dislikes           json.Number   `json:"dislikes"`
	ExtraLemmas        string        `json:"extraLemmas"`
	Md5                string        `json:"md5"`
	Views              json.Number   `json:"views"`
	Tags               []string      `json:"tags"`
	UserName           string        `json:"userName"`
	Title              string        `json:"title"`
	Description        string        `json:"description"`
	LanguageText       string        `json:"languageText"`
	LanguageCategories []string      `json:"languageCategories"`
	Subreddit          string        `json:"subreddit"`
	RedditID           string        `json:"redditId"`
	RedditIDText       string        `json:"redditIdText"`
	DomainWhitelist    []interface{} `json:"domainWhitelist"`
}

// New creates a Gfycat client.
func New(config ClientConfig) GfycatClient {
	return GfycatClient{
		ClientCredentials: ClientCredentials{
			config.ClientID,
			config.ClientSecret,
		},
	}
}

// Authenticate will authenticate using client credentials. Returns a bearer token.
func (client *GfycatClient) Authenticate() error {
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
func (client *GfycatClient) CheckToken() bool {
	if client.AccessToken == "" {
		return false
	} else {
		return true
	}
}

// GetGfycat gets a single gfycat by ID
func (client GfycatClient) GetGfycat(gfyID string) (GfycatResponse, error) {
	if ok := client.CheckToken(); !ok {
		return GfycatResponse{}, errors.New("Access token missing")
	}
	bearer := "Bearer " + client.AccessToken

	req, err := http.NewRequest(http.MethodGet, GFYCATS_URL+"/"+gfyID, nil)
	if err != nil {
		return GfycatResponse{}, err
	}
	req.Header.Add("Authorization", bearer)

	httpClient := &http.Client{}

	var resp *http.Response
	resp, err = httpClient.Do(req)
	if err != nil {
		return GfycatResponse{}, err
	}
	if resp.StatusCode == http.StatusNotFound {
		var res map[string]interface{}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return GfycatResponse{}, err
		}
		if err := json.Unmarshal(body, &res); err != nil {
			return GfycatResponse{}, err
		}
		return GfycatResponse{}, errors.New(res["errorMessage"].(string))
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return GfycatResponse{}, errors.New("unauthorized")
	}

	var gfycatResponse GfycatResponse
	err = json.NewDecoder(resp.Body).Decode(&gfycatResponse)
	if err != nil {
		return GfycatResponse{}, err
	}

	return gfycatResponse, nil
}
