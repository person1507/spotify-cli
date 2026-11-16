package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cmatthias/mapstructure"
	"github.com/jeandeaual/go-locale"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for a Spotify entity",
	Long: `Allows the user to search for something on Spotify. Searchable items include:

Song
Artist
Album
Audiobook
Episode
	`,
	Run: search,
}

type SearchResponses struct {
	Tracks     SearchResponse `json:"tracks"`
	Artists    SearchResponse `json:"artists"`
	Albums     SearchResponse `json:"albums"`
	Playlists  SearchResponse `json:"playlists"`
	Shows      SearchResponse `json:"shows"`
	Episodes   SearchResponse `json:"episodes"`
	Audiobooks SearchResponse `json:"audiobooks"`
}

type SearchResponse struct {
	Href     string  `json:"href,omitempty"`
	Limit    int     `json:"limit,omitempty"`
	Next     *string `json:"next,omitempty"`
	Offset   int     `json:"offset,omitempty"`
	Previous string  `json:"previous,omitempty"`
	Total    int64   `json:"total,omitempty"`
	Items    []any   `json:"items,omitempty"`
}

func search(cmd *cobra.Command, args []string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter your search query: ")

	var query string
	if scanner.Scan() {
		query = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		printFatal("Error reading search query: " + err.Error())
	}

	fmt.Println("")

	var types []string
	// code to get types from user
	types = []string{"track"}

	var market string
	locale, err := locale.GetLocale()
	if err != nil {
		market = "US"
	} else {
		// ex. strip "FR" from "fr-FR"
		market = locale[len(locale)-2:]
	}

	queryParams := map[string]string{
		"q":      query,
		"type":   strings.Join(types, ","),
		"market": market,
		"limit":  "10",
		"offset": "0",
	}

	var resp SearchResponses
	err = spotifyClient.Call("GET", "/search", nil, &resp, queryParams)
	if err != nil {
		printFatal("error performing search: " + err.Error())
	}

	fmt.Printf("Tracks:\n\n")
	for _, trackItem := range resp.Tracks.Items {
		var track Track
		err := mapstructure.Decode(trackItem, &track)
		if err != nil {
			printFatal("error decoding item to track: " + err.Error())
		}

		artists := []string{}
		for _, artist := range track.Artists {
			artists = append(artists, artist.Name)
		}

		fmt.Printf("%s - %s\n", track.Name, strings.Join(artists, ", "))
	}
}
