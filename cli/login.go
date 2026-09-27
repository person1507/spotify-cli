package cli

import (
	"fmt"

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
		printFatal(fmt.Sprintf("error authenticating to Spotify: %s\n", err.Error()))
	}
	fmt.Println("Success!")
}
