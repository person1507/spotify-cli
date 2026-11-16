package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

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
			"user-top-read",
			"user-follow-read",
		},
	}

	return &Client{
		BaseURL:     baseURL,
		AccountsURL: accountsURL,
		OAuthConfig: oauthConfig,
	}
}

func (c *Client) Call(method, endpoint string, body, respStruct any, queryParams map[string]string) error {
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

	q := req.URL.Query()
	for k, v := range queryParams {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	
	var reqBody []byte
	if body != nil {
		reqBody, err = io.ReadAll(req.Body)
		if err != nil {
			return err
		}
		req.Body = io.NopCloser(bytes.NewBuffer(reqBody))
	}

	if os.Getenv("DEBUG") == "1" {
		b, err := httputil.DumpRequestOut(req, false)
		if err != nil {
			return fmt.Errorf("%s", "error dumping request: "+err.Error())
		}
		fmt.Println(string(b))

		if body != nil {
			debugReqBody := io.NopCloser(bytes.NewBuffer(reqBody))
			err = prettyPrintJSON(debugReqBody)
			if err != nil {
				return fmt.Errorf("%s", "error pretty printing request body: "+err.Error())
			}
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			log.Printf("Error closing response body from call to %s %s: %s", method, endpoint, closeErr.Error())
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if os.Getenv("DEBUG") == "1" {
		b, err := httputil.DumpResponse(resp, false)
		if err != nil {
			return fmt.Errorf("%s", "error dumping response: "+err.Error())
		}
		fmt.Println(string(b))

		debugRespBody := io.NopCloser(bytes.NewBuffer(respBody))
		err = prettyPrintJSON(debugRespBody)
		if err != nil {
			return fmt.Errorf("%s", "error pretty printing response body: "+err.Error())
		}
	}

	err = json.Unmarshal(respBody, respStruct)
	if err != nil {
		return fmt.Errorf("%s", "error unmarshaling response: "+err.Error())
	}

	return nil
}

func prettyPrintJSON(body io.ReadCloser) error {
	b, _ := io.ReadAll(body)
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, b, "", "    ")
	if err != nil {
		return fmt.Errorf("%s", "error indenting body: "+err.Error())
	}

	fmt.Println(prettyJSON.String())
	return nil
}
