package cmd

import (
	"fmt"
	"github.com/jrmycanady/withings"
	"github.com/spf13/cobra"
	"log"
	"net/url"
	"strings"
)

var authGenerateRequestURLCmdVars = struct {
	redirectURL string
	scopes      string
	state       string
}{}

var authGenerateRequestURLCmd = &cobra.Command{
	Use:   "generate-request-url",
	Short: "Generates the url users must access to grant the client access to their Withings data.",
	Run: func(cmd *cobra.Command, args []string) {
		rURL, err := url.Parse(authGenerateRequestURLCmdVars.redirectURL)
		if err != nil {
			log.Fatalf("failed to parse redirect-url: %s", err)
		}

		c := withings.NewClient(ConfigOptions.ClientID, ConfigOptions.ClientSecret, *rURL)
		if ConfigOptions.Demo {
			c = withings.NewClient(ConfigOptions.ClientID, ConfigOptions.ClientSecret, *rURL, withings.WithDemoMode())
		}

		scopes := strings.Split(authGenerateRequestURLCmdVars.scopes, ",")
		for i := range scopes {
			scopes[i] = strings.TrimSpace(scopes[i])
		}

		authURL, state, err := c.GetUserAuthRequestURL(scopes, authGenerateRequestURLCmdVars.state)
		if err != nil {
			log.Fatalf("failed to generate url: %s", err)
		}

		fmt.Printf("URL: %s\n", authURL.String())
		fmt.Printf("State: %s\n", state)
	},
}

func init() {
	authGenerateRequestURLCmd.Flags().StringVar(&authGenerateRequestURLCmdVars.scopes, "scopes", "user.activity,user.metrics", "Comma separated list of scopes that will be requested for access.")
	authGenerateRequestURLCmd.Flags().StringVar(&authGenerateRequestURLCmdVars.redirectURL, "state", "", "An optional state value the will be returned by the Withings API to prevent spoofing.")
	authGenerateRequestURLCmd.MarkFlagRequired("redirect-url")

	authCmd.AddCommand(authGenerateRequestURLCmd)
}
