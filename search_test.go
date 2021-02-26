package gfycat

import (
	"encoding/json"
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGfycatClient_Search(t *testing.T) {
	mockSetup()
	defer mockTeardown()

	a := assert.New(t)

	mockedResponse := SearchResponse{
		Cursor: "cursor",
	}

	jsonRes, err := json.Marshal(mockedResponse)
	a.NoError(err)

	httpmock.RegisterResponder(http.MethodGet, SEARCH_URL, httpmock.NewStringResponder(200, string(jsonRes)))

	// Create a client
	config := ClientConfig{
		ClientID:     "client_id",
		ClientSecret: "client_secret",
	}

	// Should not have access token
	client := New(config)

	// Authenticate
	err = client.Authenticate()
	a.NoError(err)

	// Get search response
	var res SearchResponse
	res, err = client.Search("michael scott")
	fmt.Printf("Cursor: %v", res.Cursor)
	a.NoError(err)
	a.Equal(mockedResponse.Cursor, res.Cursor, "response.cursor should equal mock value")
}
