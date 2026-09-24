package spotify

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"

	"github.com/zalando/go-keyring"
)

func (c *Client) Authenticate() error {
	conf := &oauth2.Config{
		ClientID: "4b398be9bd384de1a15948c9c9dd2d3a",
		Endpoint: oauth2.Endpoint{
			AuthURL:  c.AccountsURL + "/authorize",
			TokenURL: c.AccountsURL + "/api/token",
		},
		RedirectURL: "http://127.0.0.1:8000/callback",
	}

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

	url := conf.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier))
	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)

	var token *oauth2.Token
	var err error
	select {
	case code := <-codeChan:
		token, err = conf.Exchange(context.Background(), code, oauth2.VerifierOption(verifier))

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

	err = keyring.Set("spotify-cli", "access_token", token.AccessToken)
	if err != nil {
		return fmt.Errorf("error setting access token in keyring: %s", err.Error())
	}
	err = keyring.Set("spotify-cli", "access_token_expiry", token.Expiry.String())
	if err != nil {
		return fmt.Errorf("error setting access token in keyring: %s", err.Error())
	}
	err = keyring.Set("spotify-cli", "refresh_token", token.RefreshToken)
	if err != nil {
		return fmt.Errorf("error setting access token in keyring: %s", err.Error())
	}

	return nil
}
