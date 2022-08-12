package main

import (
	"bytes"
	"fmt"
	cmd "github.com/PositionExchange/posichain-gosdk/cmd/subcommands"
	semver "github.com/hashicorp/go-version"
	"net/http"
	"os"
	"path"
	"strings"

	// Need this side effect
	_ "github.com/PositionExchange/posichain-gosdk/pkg/store"
	"github.com/spf13/cobra"
)

var (
	version     string
	commit      string
	builtAt     string
	builtBy     string
	versionLink = "https://version.posichain.org/psc"
)

func main() {
	// HACK Force usage of go implementation rather than the C based one. Do the right way, see the
	// notes one line 66,67 of https://golang.org/src/net/net.go that say can make the decision at
	// build time.
	os.Setenv("GODEBUG", "netdns=go")
	cmd.VersionWrapDump = version + "-" + commit
	cmd.RootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Show version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(os.Stderr,
				"POSICHAIN CLI (C) 2022. %v, version %v-%v (%v %v)\n",
				path.Base(os.Args[0]), version, commit, builtBy, builtAt)
			os.Exit(0)
			return nil
		},
	})
	cmd.RootCmd.AddCommand(&cobra.Command{
		Use:   "version-check",
		Short: "Check for newest version",
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, httpErr := http.Get(versionLink)
			if httpErr != nil {
				return fmt.Errorf("error when get latest version. Error: %s", httpErr)
			}
			defer resp.Body.Close()
			if resp == nil || resp.StatusCode != 200 {
				return fmt.Errorf("error response when get latest version. Http status code [%d]",
					resp.StatusCode)
			}
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			extractedCurVers := strings.Split(commit, "-")
			if len(extractedCurVers) < 1 {
				return fmt.Errorf("current version is not extracted [%s]", commit)
			}
			currentVerStr := extractedCurVers[0]
			currentVer, err := semver.NewVersion(currentVerStr)
			if err != nil {
				return fmt.Errorf("current version is invalid [%s]", currentVerStr)
			}
			latestVer, err := semver.NewVersion(buf.String())
			if err != nil {
				return fmt.Errorf("latest version is invalid [%s]", buf.String())
			}
			if currentVer.LessThan(latestVer) {
				return fmt.Errorf("warning: Using outdated version %s. Redownload to upgrade to %s",
					currentVerStr, latestVer)
			}
			fmt.Printf("Your current version %s is up-to-date\n", currentVerStr)
			return nil
		},
	})
	cmd.Execute()
}
