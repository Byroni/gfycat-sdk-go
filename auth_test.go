package gfycat

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAuthClient(t *testing.T) {
	a := assert.New(t)

	clientID := "client_id"
	clientSecret := "client_secret"
	var grantType string

	authClient := NewAuthClient(clientID, clientSecret)

	a.Equal(authClient.GrantType, grantType, "GrantType should be " + grantType)
	a.Equal(authClient.ClientID, clientID, "ClientID should be " + clientID)
	a.Equal(authClient.ClientSecret, clientSecret, "ClientSecret should be " + clientSecret)
}

func TestAuthRequest_GetAccessToken(t *testing.T) {
	a := assert.New(t)

	mockedResponse := `
		{
			"token_type":"bearer",
			"scope":"",
			"expires_in":3600,
			"access_token":"access_token"
		}
	`

	httpmock.Activate()
	httpmock.RegisterResponder("POST", AUTH_TOKEN_URL, httpmock.NewStringResponder(200, mockedResponse))

	authClient := NewAuthClient("client_credentials", "client_secret")

	authResponse, err := authClient.GetAccessToken()
	if err != nil {
		t.Errorf("Failed with error %s", err.Error())
		return
	}

	httpmock.DeactivateAndReset()

	a.NoError(err, "Error should be nil")
	a.Equal(authResponse.AccessToken, "access_token", "Access token should be equal to the mocked response")
	a.Empty(authResponse.RefreshToken, "Refresh token should be empty")
}

func TestAuthRequest_GetAccessTokenFail(t *testing.T) {
	a := assert.New(t)

	httpmock.Activate()
	httpmock.RegisterResponder("POST", AUTH_TOKEN_URL, httpmock.NewStringResponder(401, ""))

	authClient := NewAuthClient("client_credentials", "client_secret")

	httpmock.DeactivateAndReset()

	_, err := authClient.GetAccessToken()
	a.Error(err, "Error should not be nil")
	a.Contains(err.Error(), "Request failed with status code")
}

func TestAuthRequest_RefreshAccessToken(t *testing.T) {
	a := assert.New(t)

	mockedResponse := `
		{
			"token_type":"bearer",
			"refresh_token_expires_in":5184000,
			"refresh_token":"refresh_token",
			"scope":"",
			"resource_owner":"username",
			"expires_in":3600,
			"access_token":"access_token"}
	`

	httpmock.Activate()
	httpmock.RegisterResponder("POST", AUTH_TOKEN_URL, httpmock.NewStringResponder(200, mockedResponse))

	authClient := NewAuthClient("client_credentials", "client_secret")

	authResponse, err := authClient.RefreshAccessToken("access_token")
	if err != nil {
		t.Errorf("Failed with error %s", err.Error())
		return
	}

	httpmock.DeactivateAndReset()

	a.NoError(err, "Error should be nil")
	a.Equal(authResponse.AccessToken, "access_token", "Access token should be equal to the mocked response")
	a.Equal(authResponse.RefreshToken, "refresh_token", "Refresh token should be equal to the mocked response")
}