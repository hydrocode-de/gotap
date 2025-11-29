package cmd

import (
	"fmt"

	"github.com/hydrocode-de/gotap/internal/io"
	"github.com/hydrocode-de/gotap/internal/validation"
	"github.com/spf13/cobra"
)

var verbose bool

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify tool-spec metadata",
	Long: `Verify the tool-spec metadata against the tool-spec schema.

This command will verify the tool-spec metadata against the tool-spec schema.
It will collect and return all verification errors.
`,
	Run: verify,
}

func verify(cmd *cobra.Command, args []string) {
	// run validation
	validation, err := validation.LoadAndValidateSpec(args)
	cobra.CheckErr(err)

	errorCount := validation.ErrorCount()
	warningCount := validation.WarningCount()
	hasErrors := errorCount > 0
	hasWarnings := warningCount > 0

	warnings := validation.Warnings
	errors := validation.Errors

	if !hasErrors && !hasWarnings {
		fmt.Println("OK")
		return
	} else if hasErrors {
		fmt.Println("FAIL")
	} else {
		fmt.Println("WARN")
	}

	for _, warning := range warnings {
		if verbose {
			fmt.Println("---")
		}
		fmt.Println(io.WriteValidationError(warning, verbose))
	}

	for _, err := range errors {
		if verbose {
			fmt.Println("---")
		}
		fmt.Println(io.WriteValidationError(err, verbose))
	}

	if verbose && (hasErrors || hasWarnings) {
		fmt.Println("--------------------------------")
		fmt.Printf("ERRORS: %d     WARNINGS: %d\n", errorCount, warningCount)
		fmt.Println("--------------------------------")
	}
}

func init() {
	verifyCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	rootCmd.AddCommand(verifyCmd)
}
