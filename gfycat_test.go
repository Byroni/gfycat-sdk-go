package gfycat

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

var mockGfycatID = "mockID"

/**
Run whole test file to setup and teardown mocks
*/
func mockSetup() {
	httpmock.Activate()
	mockedAuthResponse := `
		{
			"token_type":"bearer",
			"scope":"",
			"expires_in":3600,
			"access_token":"access_token"
		}
	`
	httpmock.RegisterResponder("POST", AUTH_TOKEN_URL, httpmock.NewStringResponder(200, mockedAuthResponse))
}

func TestGfycatClient_CheckToken(t *testing.T) {
	mockSetup()
	a := assert.New(t)

	config := ClientConfig{
		ClientID:     "client_id",
		ClientSecret: "client_secret",
	}

	// Should not have access token
	client := New(config)

	// do not authenticate

	ok := client.CheckToken()
	a.False(ok)

	// authenticate now
	err := client.Authenticate()
	a.NoError(err)

	ok = client.CheckToken()
	a.True(ok)

	a.NotEmpty(client.AccessToken)

	mockTeardown()
}

func TestGetGfycat(t *testing.T) {
	mockSetup()
	a := assert.New(t)

	mockedResponse := `
		{
			"gfyItem": {
				"gfyId": "id",
				"likes": 0
			}
		}
	`

	httpmock.RegisterResponder("GET", GFYCATS_URL+"/"+mockGfycatID, httpmock.NewStringResponder(200, mockedResponse))

	config := ClientConfig{
		ClientID:     "client_id",
		ClientSecret: "client_secret",
	}

	client := New(config)

	// Check for access token
	_, err := client.GetGfycat(mockGfycatID)
	a.Error(err)
	a.EqualError(err, "Access token missing")

	// Authenticate now
	err = client.Authenticate()
	a.NoError(err)

	gfycatResponse, err := client.GetGfycat(mockGfycatID)
	fmt.Println(err)
	a.NoError(err)
	a.NotEmpty(gfycatResponse.GfyItem.GfyID)
	a.Equal("id", gfycatResponse.GfyItem.GfyID, "GfyID must equal mocked value")
	a.Equal(0, gfycatResponse.GfyItem.Likes, "Likes must equal mocked value")

	mockTeardown()
}

func TestClient_GetGfycat_successfully_handles_404(t *testing.T) {
	mockSetup()
	defer mockTeardown()

	mockedResponse := `
		{
			"errorMessage": "does not exist."
		}
	`
	httpmock.RegisterResponder("GET", GFYCATS_URL+"/"+mockGfycatID, httpmock.NewStringResponder(404, mockedResponse))

	a := assert.New(t)

	client := New(ClientConfig{})

	err := client.Authenticate()
	a.NoError(err)

	_, err = client.GetGfycat(mockGfycatID)
	a.Error(err)
	a.Contains(err.Error(), "does not exist.")

}

func mockTeardown() {
	httpmock.DeactivateAndReset()
}
