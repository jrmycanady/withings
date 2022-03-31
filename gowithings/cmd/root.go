package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var ConfigOptions = struct {
	ClientID                    string
	ClientSecret                string
	SkipCertificateVerification bool
	Demo                        bool
}{}

var rootCmd = &cobra.Command{
	Use: "gowithings",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&ConfigOptions.ClientID, "client-id", "i", "", "The client id provided by Withings for the client.")
	rootCmd.PersistentFlags().StringVarP(&ConfigOptions.ClientSecret, "client-secret", "s", "", "The client secret provided by Withings for the client.")
	rootCmd.PersistentFlags().StringVarP(&authGenerateRequestURLCmdVars.redirectURL, "redirect-url", "u", "", "The URL the Withings API should redirect back to.")

	rootCmd.PersistentFlags().BoolVar(&ConfigOptions.SkipCertificateVerification, "skip-certificate-verification", false, "The client secret provided by Withings for the client.")
	rootCmd.PersistentFlags().BoolVar(&ConfigOptions.Demo, "demo-mode", false, "Denotes if all API calls should use demo mode.")
}
