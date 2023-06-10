package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"code.nkcmr.net/datever"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func rootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "datever",
	}
	cmd.AddCommand(incrementCommand())
	cmd.AddCommand(parseCommand())
	return cmd
}

func parseCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "parse VERSION",
		Args: cobra.ExactArgs(1),
		Run: runWithError(func(c *cobra.Command, s []string) error {
			verText := s[0]
			ver, err := datever.Parse(verText)
			if err != nil {
				return errors.Wrap(err, "failed to parse input version")
			}
			jsbytes, err := json.Marshal(ver)
			if err != nil {
				return errors.Wrap(err, "failed to json encode version")
			}
			_, _ = os.Stdout.Write(jsbytes)
			return nil
		}),
	}
}

func incrementCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "increment CURRENT_VERSION",
		Args: cobra.ExactArgs(1),
		Run: runWithError(func(c *cobra.Command, s []string) error {
			verText := s[0]
			ver, err := datever.Parse(verText)
			if err != nil {
				return errors.Wrap(err, "failed to parse input version")
			}
			nextVer := ver.Increment(time.Now())
			fmt.Print(nextVer.String())
			return nil
		}),
	}
}

func runWithError(rf func(*cobra.Command, []string) error) func(*cobra.Command, []string) {
	return func(c *cobra.Command, s []string) {
		if err := rf(c, s); err != nil {
			fmt.Printf("error: %s\n", err.Error())
			os.Exit(1)
		}
	}
}

func main() {
	rootCommand().Execute()
}
