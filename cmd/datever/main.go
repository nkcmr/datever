package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"code.nkcmr.net/datever"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

func trimLineSpace(s string) string {
	lines := strings.Split(s, "\n")
	outlines := make([]string, 0, len(lines))
	for i := range lines {
		prevLineString := len(outlines) > i && outlines[i-1] == "\n"
		trimmedLine := strings.TrimSpace(lines[i])
		if trimmedLine == "" && (prevLineString || len(outlines) == 0) {
			continue
		}
		outlines = append(outlines, trimmedLine)
	}
	return strings.Join(outlines, "\n")
}

func rootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "datever",
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		Long: trimLineSpace(`
			A tool to deal with 'date-based versioning'.
			
			For software that does not need to convey breakages with other software (like
			a library would), date versioning is a good way to label releases in a way that
			allows humans to understand the rough age of a release and do easy comparing.

			This tool enables the usage of date version strings and helps integrate them
			with release scripts.
		`),
	}
	cmd.AddCommand(incrementCommand())
	cmd.AddCommand(parseCommand())
	cmd.AddCommand(compareCommand())
	return cmd
}

func parseCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "parse VERSION",
		Short: "parse a datever string (2024.2.9) and return a json breakdown of its parts",
		Args:  cobra.ExactArgs(1),
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

func operatorGreaterThan(a, b datever.Version) bool {
	// a > b
	return datever.Compare(a, b) == 1
}

func operatorEqual(a, b datever.Version) bool {
	return datever.Compare(a, b) == 0
}

func operatorGreaterThanOrEqual(a, b datever.Version) bool {
	return operatorGreaterThan(a, b) || operatorEqual(a, b)
}

func operatorLessThan(a, b datever.Version) bool {
	return datever.Compare(a, b) == -1
}

func operatorLessThanOrEqual(a, b datever.Version) bool {
	return operatorLessThan(a, b) || operatorEqual(a, b)
}

var operatorMappings = map[string]func(a, b datever.Version) bool{
	">":   operatorGreaterThan,
	"gt":  operatorGreaterThan,
	">=":  operatorGreaterThanOrEqual,
	"gte": operatorGreaterThanOrEqual,
	"=":   operatorEqual,
	"==":  operatorEqual,
	"eq":  operatorEqual,
	"<":   operatorLessThan,
	"lt":  operatorLessThan,
	"<=":  operatorLessThanOrEqual,
	"lte": operatorLessThanOrEqual,
}

func compareCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "compare VERSION_A OPERATOR VERSION_B",
		Short: fmt.Sprintf("compare 2 versions (available operators: %s)", strings.Join(maps.Keys(operatorMappings), ", ")),
		Args:  cobra.ExactArgs(3),
		Run: runWithError(func(c *cobra.Command, s []string) error {
			va, err := datever.Parse(s[0])
			if err != nil {
				return errors.Wrap(err, "failed to parse LHS version")
			}
			vb, err := datever.Parse(s[2])
			if err != nil {
				return errors.Wrap(err, "failed to parse RHS version")
			}
			opfunc, ok := operatorMappings[strings.ToLower(s[1])]
			if !ok {
				return fmt.Errorf("unknown operator '%s'", s[1])
			}
			if !opfunc(va, vb) {
				fmt.Println("false")
				os.Exit(1)
			} else {
				fmt.Println("true")
			}
			return nil
		}),
	}
}

func incrementCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "increment CURRENT_VERSION",
		Short: "take a date version and increment it (2000.1.0 => 2025.1.0, 2024.3.4 => 2024.3.5)",
		Args:  cobra.ExactArgs(1),
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
