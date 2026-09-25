package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var userInfoCmd = &cobra.Command{
	Use:   "user-info",
	Short: "Get info about the currently logged-in user.",
	Long:  "Get info about the currently logged-in user.",
	Run:   getUserInfo,
}

type UserProfile struct {
	AccountID       string `json:"account_id"`
	Country         string `json:"country"`
	DisplayName     string `json:"display_name"`
	Email           string `json:"email"`
	ExplicitContent struct {
		FilterEnabled bool `json:"filter_enabled"`
		FilterLocked  bool `json:"filter_locked"`
	} `json:"explicit_content"`
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Followers struct {
		Total int `json:"total"`
	} `json:"followers"`
	Href    string `json:"href"`
	ID      string `json:"id"`
	Product string `json:"product"`
	Type    string `json:"type"`
	URI     string `json:"uri"`
}

func getUserInfo(cmd *cobra.Command, args []string) {
	resp, err := spotifyClient.Call("GET", "/me", nil)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}

	var profile UserProfile
	err = json.Unmarshal(resp, &profile)
	if err != nil {
		fmt.Println("error unmarshaling response: " + err.Error())
		fmt.Println("Raw body response:\n" + string(resp))
		os.Exit(1)
	}

	fmt.Println("Name: " + profile.DisplayName)
	fmt.Println("Username: " + profile.ID)
	fmt.Println("Email: " + profile.Email)
	fmt.Println("Tier: " + profile.Product)
	fmt.Printf("Followers: %d", profile.Followers.Total)
}
