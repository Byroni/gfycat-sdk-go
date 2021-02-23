package gfycat

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
	client, err := New(config)
	a.NoError(err)

	// do not authenticate

	ok := client.CheckToken()
	a.False(ok)

	// authenticate now
	err = client.Authenticate()
	a.NoError(err)

	ok = client.CheckToken()
	a.True(ok)

	mockTeardown()
}

func TestGetGfycat(t *testing.T) {
	mockSetup()
	a := assert.New(t)
	gfycatID :=  "mockID"

	mockedResponse := `
		{
			"GfyItem": {
				"GfyID": "id"
			}
		}		
	`

	httpmock.RegisterResponder("GET", GFYCATS_URL+"/"+gfycatID, httpmock.NewStringResponder(200, mockedResponse))

	config := ClientConfig{
		ClientID:     "client_id",
		ClientSecret: "client_secret",
	}

	client, err := New(config)
	a.NoError(err)

	// Check for access token
	_, err = client.GetGfycat(gfycatID)
	a.Error(err)
	a.EqualError(err, "Access token missing")

	// Authenticate now
	err = client.Authenticate()
	a.NoError(err)

	gfycatResponse, err := client.GetGfycat(gfycatID)

	a.NoError(err)
	a.NotEmpty(gfycatResponse.GfyItem.GfyID)
	a.Equal("id", gfycatResponse.GfyItem.GfyID, "GfyID must equal mocked value")

	mockTeardown()
}

func mockTeardown() {
	httpmock.DeactivateAndReset()
}