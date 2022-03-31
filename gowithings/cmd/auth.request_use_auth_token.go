package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/jrmycanady/withings"
	"github.com/spf13/cobra"
	"log"
	"net/url"
)

var authCode string

var authRequestUserAuthToken = &cobra.Command{
	Use:   "request-user-auth-token",
	Short: "Requests a user auth token and provides the results.",
	Run: func(cmd *cobra.Command, args []string) {
		rURL, err := url.Parse(authGenerateRequestURLCmdVars.redirectURL)
		if err != nil {
			log.Fatalf("failed to parse redirect-url: %s", err)
		}

		c := withings.NewClient(ConfigOptions.ClientID, ConfigOptions.ClientSecret, *rURL)
		resp, err := c.GetUserAccessToken(authCode)
		if err != nil {
			log.Fatalf("Failed to get token: %s", err)
		}

		if resp.Status != 0 {
			log.Fatalf("Failed to get token with status response %d", resp.Status)
		}

		token, err := json.MarshalIndent(resp.AccessToken, "", " ")
		if err != nil {
			fmt.Errorf("Failed to mashal response: %s", err)
		}

		log.Println(string(token))
	},
}

func init() {
	authRequestUserAuthToken.Flags().StringVar(&authCode, "code", "", "The authentication code for the user.")
	authRequestUserAuthToken.MarkFlagRequired("code")

	authCmd.AddCommand(authRequestUserAuthToken)
}
