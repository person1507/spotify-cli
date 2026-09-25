package spotify

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/zalando/go-keyring"
)

func (c *Client) Authenticate() error {
	verifier := oauth2.GenerateVerifier()
	state := rand.Text()

	codeChan := make(chan string, 1)
	errChan := make(chan string, 1)
	mux := http.NewServeMux()

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		err := r.URL.Query().Get("error")
		stateReturned := r.URL.Query().Get("state")

		if stateReturned != state {
			errChan <- "returned state did not match state passed in"
		}

		if code == "" && err == "" {
			http.Error(w, "Missing code", http.StatusBadRequest)
			return
		} else if code != "" {
			codeChan <- code
		} else {
			errChan <- err
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		html := `<!DOCTYPE html>
		<html>
		<head><title>Go HTML</title></head>
		<body>
			<h2>Callback received, you may now close this window.</h2>
		</body>
		</html>`
		
		fmt.Fprint(w, html)
	})

	server := &http.Server{
		Addr:    "127.0.0.1:8000",
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("callback server could not start: %s", err.Error())
		}
	}()

	url := c.OAuthConfig.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier))
	fmt.Printf("Visit the URL for the auth dialog: %v\n\n", url)

	var token *oauth2.Token
	var err error
	select {
	case code := <-codeChan:
		token, err = c.OAuthConfig.Exchange(context.Background(), code, oauth2.VerifierOption(verifier))

	case authErr := <-errChan:
		return fmt.Errorf("authorization request failed: %s", authErr)

	case <-time.After(2 * time.Minute):
		server.Shutdown(context.Background())
		return fmt.Errorf("timed out waiting for Spotify authorization")
	}

	server.Shutdown(context.Background())

	if err != nil {
		return fmt.Errorf("error during exchange: %s", err.Error())
	}

	err = setAccessTokenAndExpiry(token)
	if err != nil {
		return err
	}

	err = keyring.Set("spotify-cli", "refresh_token", token.RefreshToken)
	if err != nil {
		return fmt.Errorf("error setting refresh token in keyring: %s", err.Error())
	}

	return nil
}

func (c *Client) getAccessToken() (string, error) {
	accessToken, err := keyring.Get("spotify-cli", "access_token")
	if err != nil {
		return "", fmt.Errorf("failed to get access token from keyring: %s", err.Error())
	}

	expiry, err := keyring.Get("spotify-cli", "access_token_expiry")
	if err != nil {
		return "", fmt.Errorf("failed to get token expiry from keyring: %s", err.Error())
	}

	refreshToken, err := keyring.Get("spotify-cli", "refresh_token")
	if err != nil {
		return "", fmt.Errorf("failed to get access token in keyring: %s", err.Error())
	}

	cleanStr := strings.Split(expiry, " m=")[0]
	layout := "2006-01-02 15:04:05.000000 -0700 MST"

	t, err := time.Parse(layout, cleanStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse token expiry: %s", err.Error())
	}

	if !t.After(time.Now()) {
		tokenSource := c.OAuthConfig.TokenSource(context.Background(), &oauth2.Token{
			RefreshToken: refreshToken,
		})

		token, err := tokenSource.Token()
		if err != nil {
			return "", fmt.Errorf("failed to refresh token: %s", err.Error())
		}

		err = setAccessTokenAndExpiry(token)
		if err != nil {
			return "", err
		}
		accessToken = token.AccessToken
	}

	return accessToken, nil
}

func setAccessTokenAndExpiry(token *oauth2.Token) error {
	err := keyring.Set("spotify-cli", "access_token", token.AccessToken)
	if err != nil {
		return fmt.Errorf("error setting access token in keyring: %s", err.Error())
	}

	err = keyring.Set("spotify-cli", "access_token_expiry", token.Expiry.String())
	if err != nil {
		return fmt.Errorf("error setting access token expiry token in keyring: %s", err.Error())
	}

	return nil
}
