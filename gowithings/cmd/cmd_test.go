package cmd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"withings"
)

func TestExecute_Root(t *testing.T) {

	tests := map[string]struct {
		args   []string
		stdErr string
		stdOut string
	}{
		"gowithings": {
			args:   []string{""},
			stdErr: "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var stdErrBuff bytes.Buffer
			var stdOutBuff bytes.Buffer

			rootCmd.SetArgs(test.args)
			rootCmd.SetErr(&stdErrBuff)
			rootCmd.SetOut(&stdOutBuff)

			rootCmd.Execute()

			switch test.stdErr {
			case "":
				assert.Empty(t, stdErrBuff.String())
				assert.NotEmpty(t, stdOutBuff.String())
			default:
				assert.NotEmpty(t, stdErrBuff.String())
			}

		})
	}
}

func TestExecute_Auth_GenerateRequestURL(t *testing.T) {

	tests := map[string]struct {
		redirectURL string
		state       string
		scopes      string
		stdErr      string
		stdOut      string
	}{
		"With URL, state and scopes": {
			redirectURL: "example.com/",
			state:       "TESTSTATE",
			scopes:      strings.Join([]string{withings.ScopeUserActivity, withings.ScopeUserMetrics}, ","),
			stdErr:      "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var stdErrBuff bytes.Buffer
			var stdOutBuff bytes.Buffer

			rootCmd.SetArgs([]string{})
			rootCmd.SetErr(&stdErrBuff)
			rootCmd.SetOut(&stdOutBuff)

			rootCmd.Execute()

			switch test.stdErr {
			case "":
				assert.Empty(t, stdErrBuff.String())
				assert.NotEmpty(t, stdOutBuff.String())
				assert.Contains(t, stdOutBuff.String(), "STATEVALUE")
			default:
				assert.NotEmpty(t, stdErrBuff.String())
			}

		})
	}
}
