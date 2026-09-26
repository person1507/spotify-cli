package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

type Client struct {
	BaseURL     string
	AccountsURL string
	OAuthConfig *oauth2.Config
}

func New() *Client {
	baseURL := "https://api.spotify.com/v1"
	accountsURL := "https://accounts.spotify.com"

	oauthConfig := &oauth2.Config{
		ClientID: "4b398be9bd384de1a15948c9c9dd2d3a",
		Endpoint: oauth2.Endpoint{
			AuthURL:  accountsURL + "/authorize",
			TokenURL: accountsURL + "/api/token",
		},
		RedirectURL: "http://127.0.0.1:8000/callback",
		Scopes: []string{
			"user-read-private",
			"user-read-email",
		},
	}

	return &Client{
		BaseURL:     baseURL,
		AccountsURL: accountsURL,
		OAuthConfig: oauthConfig,
	}
}

func (c *Client) Call(method, endpoint string, body, respStruct any) error {
	token, err := c.getAccessToken()
	if err != nil {
		return err
	}

	var bodyReader io.Reader
	if body != nil {
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(bodyJSON)
	}

	client := http.Client{}

	req, err := http.NewRequest(method, c.BaseURL+endpoint, bodyReader)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respBody, respStruct)
	if err != nil {
		return fmt.Errorf("%s", "error unmarshaling response: " + err.Error())
	}

	return nil
}
