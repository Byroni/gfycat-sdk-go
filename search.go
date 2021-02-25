package gfycat

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// SearchRequest describes the query parameters for the Gfycat search API
type SearchRequest struct {
	SearchText string `json:"search_text"`
	Count      int    `json:"count,omitempty"`
	Cursor     int    `json:"cursor,omitempty"`
}

// SearchResponse describes the search API response
type SearchResponse struct {
	Cursor  string `json:"cursor"`
	Gfycats []struct {
		GfyItem
		Tags []string `json:"tags"`
		UserData  struct {
			Name            string `json:"name"`
			ProfileImageURL string `json:"profileImageUrl"`
			URL             string `json:"url"`
			Username        string `json:"username"`
			Followers       int    `json:"followers"`
			Subscription    int    `json:"subscription"`
			Following       int    `json:"following"`
			ProfileURL      string `json:"profileUrl"`
			Views           int    `json:"views"`
			Verified        bool   `json:"verified"`
		} `json:"userData"`
		ContentUrls         struct {
			Max2MbGif struct {
				URL    string `json:"url"`
				Size   int    `json:"size"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"max2mbGif"`
			Webp struct {
				URL    string `json:"url"`
				Size   int    `json:"size"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"webp"`
			Max1MbGif struct {
				URL    string `json:"url"`
				Size   int    `json:"size"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"max1mbGif"`
			One00PxGif struct {
				URL    string `json:"url"`
				Size   int    `json:"size"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"100pxGif"`
			MobilePoster struct {
				URL    string `json:"url"`
				Size   int    `json:"size"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"mobilePoster"`
			Mp4 struct {
				URL    string `json:"url"`
				Size   int    `json:"size"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"mp4"`
			Webm struct {
				URL    string `json:"url"`
				Size   int    `json:"size"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"webm"`
			Max5MbGif struct {
				URL    string `json:"url"`
				Size   int    `json:"size"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"max5mbGif"`
			LargeGif struct {
				URL    string `json:"url"`
				Size   int    `json:"size"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"largeGif"`
			Mobile struct {
				URL    string `json:"url"`
				Size   int    `json:"size"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"mobile"`
		} `json:"content_urls"`
	} `json:"gfycats"`
	Related []string `json:"related"`
	Found   int      `json:"found"`
}

// Search takes a search string and queries against the Gfycat library. Gfycat site search is limited to
func (client GfycatClient) Search(searchTerm string) (SearchResponse, error) {
	// TODO: Add cursor and count support
	if ok := client.CheckToken(); !ok {
		return SearchResponse{}, errors.New("Access token missing")
	}
	bearer := "Bearer " + client.AccessToken

	req, err := http.NewRequest("GET", SEARCH_URL, nil)
	if err != nil {
		return SearchResponse{}, err
	}
	// Add authorization header
	req.Header.Add("Authorization", bearer)

	fmt.Println(bearer)
	// Add query params
	q := req.URL.Query()
	q.Add("search_text", searchTerm)
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())
	httpClient := &http.Client{}

	var resp *http.Response
	resp, err = httpClient.Do(req)
	if err != nil {
		return SearchResponse{}, err
	}
	if resp.StatusCode == 403 {
		fmt.Printf("%+v", resp)
		return SearchResponse{}, errors.New("unauthorized")
	}

	var searchResponse SearchResponse
	err = json.NewDecoder(resp.Body).Decode(&searchResponse)
	if err != nil {
		return SearchResponse{}, err
	}

	fmt.Printf("%+v", searchResponse)

	return SearchResponse{}, nil
}
