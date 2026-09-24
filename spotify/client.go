package spotify

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Client struct {
	BaseURL     string
	AccountsURL string
}

func New() *Client {
	return &Client{
		BaseURL:     "https://api.spotify.com/v1",
		AccountsURL: "https://accounts.spotify.com",
	}
}

func (c *Client) CallFullURL(method, endpoint string, body any) ([]byte, error) {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	client := http.Client{}

	req, err := http.NewRequest(method, c.BaseURL+endpoint, bytes.NewReader(bodyJSON))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
