package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log into Spotify",
	Long:  "This logs the user into Spotify via OAuth with PKCE.",
	Run:   login,
}

func login(cmd *cobra.Command, args []string) {
	err := spotifyClient.Authenticate()
	if err != nil {
		fmt.Printf("Error authenticating to Spotify: %s\n", err.Error())
		os.Exit(1)
	} else {
		fmt.Println("Success!")
	}
}
