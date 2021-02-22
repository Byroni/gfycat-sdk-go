package gfycat

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

const GRANT_TYPE_CLIENT_CREDENTIALS = "client_credentials"
/**
	These are currently unused until the other authentication methods are implemented.
 */
const GRANT_TYPE_PASSWORD = "password"
const GRANT_TYPE_REFRESH = "refresh"
const GRANT_TYPE_AUTH_CODE = "authorization_code"
const GRANT_TYPE_CONVERT_CODE = "convert_code"
const GRANT_TYPE_CONVERT_TOKEN = "convert_token"
const GRANT_TYPE_REQUEST_TOKEN = "request_token"
const GRANT_TYPE_CONVERT_REQUEST_TOKEN = "convert_request_token"
const GRANT_TYPE_PROVIDER_TOKEN = "provider_token"

// AuthResponse contains Gfycat's authenticated bearer token. The token defaults to 1 hour expiry.
type AuthResponse struct {
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
}

// AuthRequest contains client credentials and grant type
type AuthRequest struct {
	GrantType string `json:"grant_type"`
	ClientID string `json:"client_id"`
	ClientSecret  string `json:"client_secret"`
}

// NewAuthClient creates a new client used for authenticating against Gfycat API.
// It returns an instance of AuthRequest.
func NewAuthClient(clientID string, clientSecret string) AuthRequest {
	return AuthRequest{
		ClientID: clientID,
		ClientSecret: clientSecret,
	}
}

// GetAccessToken will attempt to get a new access token from Gfycat.
// If successful, it returns an instance of AuthResponse.
// If unsuccessful, it returns an error
func (authRequest AuthRequest) GetAccessToken() (AuthResponse, error) {
	authRequest.GrantType = GRANT_TYPE_CLIENT_CREDENTIALS

	payload, err := json.Marshal(authRequest)
	if err != nil {
		return AuthResponse{}, err
	}

	var resp *http.Response
	resp, err = http.Post(AUTH_TOKEN_URL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return AuthResponse{}, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return AuthResponse{}, errors.New("Request failed with status code" + string(rune(resp.StatusCode)))
	}

	var authResponse AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	if err != nil {
		return AuthResponse{}, err
	}
	return authResponse, nil
}

// TODO: Refresh token

// TODO: Password based authentication

// TODO: Browser based authentication