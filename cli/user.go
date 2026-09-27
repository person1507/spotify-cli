package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

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
	ExternalURLs ExternalURL `json:"external_urls"`
	Followers    struct {
		Total int64 `json:"total"`
	} `json:"followers"`
	Href    string `json:"href"`
	ID      string `json:"id"`
	Product string `json:"product"`
	Type    string `json:"type"`
	URI     string `json:"uri"`
}

func getUserInfo(cmd *cobra.Command, args []string) {
	var profile UserProfile
	err := spotifyClient.Call("GET", "/me", nil, &profile, nil)
	if err != nil {
		printFatal("Error: " + err.Error())
	}

	fmt.Println("Name: " + profile.DisplayName)
	fmt.Println("Username: " + profile.ID)
	fmt.Println("Email: " + profile.Email)
	fmt.Println("Tier: " + profile.Product)
	fmt.Printf("Followers: %d\n", profile.Followers.Total)
}

var topItemsCmd = &cobra.Command{
	Use:   "top-items [artists|songs]",
	Short: "Get user's top items",
	Long:  "Get user's top artists and songs",
	Args:  cobra.ExactArgs(1),
	Run:   getUserTopItems,
}

type topItemsResponse struct {
	Href     string  `json:"href"`
	Limit    int64   `json:"limit"`
	Next     *string `json:"next"`
	Offset   int64   `json:"offset"`
	Previous *string `json:"previous"`
	Total    int64   `json:"total"`
	Items    any     `json:"items"`
}

func getUserTopItems(cmd *cobra.Command, args []string) {
	entity := args[0]
	if entity != "artists" && entity != "songs" {
		printFatal("error: argument must be either 'artists' or 'songs'")
		os.Exit(1)
	}
	if entity == "songs" {
		entity = "tracks"
	}

	queryParams := map[string]string{
		"time_range": "short_term",
		"limit":      "10",
		"offset":     "0",
	}

	var topItems topItemsResponse
	err := spotifyClient.Call("GET", "/me/top/"+entity, nil, &topItems, queryParams)
	if err != nil {
		printFatal(err.Error())
	}

	if entity == "artists" {
		itemsJSON, _ := json.Marshal(topItems.Items)
		var artists []Artist
		err := json.Unmarshal(itemsJSON, &artists)
		if err != nil {
			printFatal("error: type of items returned is not artists: " + err.Error())
		}
		fmt.Println("Your Top 10 Artists:")
		for _, artist := range artists {
			fmt.Println(artist.Name)
		}
	} else {
		itemsJSON, _ := json.Marshal(topItems.Items)
		var tracks []Track
		err := json.Unmarshal(itemsJSON, &tracks)
		if err != nil {
			printFatal("error: type of items returned is not tracks: " + err.Error())
		}
		fmt.Println("Your top 10 songs:")
		for _, track := range tracks {
			var artistList strings.Builder
			for i, artist := range track.Artists {
				if i == 0 {
					artistList.WriteString(artist.Name)
				} else {
					artistList.WriteString(", ")
					artistList.WriteString(artist.Name)
				}
			}
			fmt.Println(track.Name + " - " + artistList.String())
		}
	}
}

var followedArtistsCmd = &cobra.Command{
	Use:   "followed-artists",
	Short: "Get the user's list of followed artists",
	Long:  "Get the user's list of followed artists",
	Run:   getFollowedArtists,
}

type FollowedArtistResponse struct {
	Artists struct {
		Href    string  `json:"href"`
		Limit   int64   `json:"limit"`
		Next    *string `json:"next"`
		Cursors struct {
			Before string `json:"before"`
			After  string `json:"after"`
		} `json:"cursors"`
		Total int64    `json:"total"`
		Items []Artist `json:"items"`
	} `json:"artists"`
}

func getFollowedArtists(cmd *cobra.Command, args []string) {
	queryParams := map[string]string{
		"type": "artist",
	}

	var followedArtists FollowedArtistResponse
	err := spotifyClient.Call("GET", "/me/following", nil, &followedArtists, queryParams)
	if err != nil {
		printFatal(err.Error())
	}

	var artists []Artist
	artists = append(artists, followedArtists.Artists.Items...)
	for followedArtists.Artists.Next != nil {
		queryParams["after"] = followedArtists.Artists.Cursors.After
		err := spotifyClient.Call("GET", "/me/following", nil, &followedArtists, queryParams)
		if err != nil {
			printFatal(err.Error())
		}
		artists = append(artists, followedArtists.Artists.Items...)
	}

	fmt.Printf("Your followed artists (count %d):\n", followedArtists.Artists.Total)
	for _, artist := range artists {
		fmt.Println(artist.Name)
	}
}
