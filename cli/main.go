package cli

import (
	"log"
	"os"

	"github.com/person1507/spotify-cli/spotify"
	"github.com/spf13/cobra"
)

var spotifyClient *spotify.Client

var rootCmd = &cobra.Command{
	Use:   "spotify",
	Short: "A CLI for interacting with Spotify",
	Long:  "A CLI for interacting with Spotify",
}

func init() {
	rootCmd.AddCommand(loginCmd, userInfoCmd)
}

func Main() {
	spotifyClient = spotify.New()

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
