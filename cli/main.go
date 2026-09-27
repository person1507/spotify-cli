package cli

import (
	"log"
	"os"

	"github.com/fatih/color"
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
	rootCmd.AddCommand(loginCmd, userInfoCmd, topItemsCmd, followedArtistsCmd)
}

func printFatal(err string) {
	boldRed := color.New(color.FgRed, color.Bold)
	boldRed.Println(err)
	os.Exit(1)
}

func Main() {
	spotifyClient = spotify.New()

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
