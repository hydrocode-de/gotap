/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/hydrocode-de/gotap/internal/config"
	"github.com/hydrocode-de/gotap/internal/io"
	"github.com/hydrocode-de/tool-spec-go/validate"
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
	v := config.GetViper()
	specFile := v.GetString("spec_file")
	inputFile := v.GetString("input_file")
	citationFile := v.GetString("citation_file")
	licenseFile := v.GetString("license_file")

	spec, err := io.ReadSpecFile(specFile)
	cobra.CheckErr(err)

	input, err := io.ReadInputFile(inputFile)
	cobra.CheckErr(err)

	toolname, err := config.ResolveToolname(args, input)
	cobra.CheckErr(err)

	toolSpec, err := spec.GetTool(toolname)
	cobra.CheckErr(err)

	warnings := make([]validate.ValidationError, 0)
	citation, err := io.ReadCitationFile(citationFile)
	if err != nil {
		warnings = append(warnings, validate.ValidationError{
			Field:    "/src/CITATION.cff",
			Name:     "CITATION.cff",
			Message:  "No citation file found. We recommend adding one.",
			Type:     validate.ErrorType("WARNING"),
			Expected: "CITATION.cff",
			Actual:   "None",
		})
	} else {
		toolSpec.Citation = citation
	}
	toolInput, err := input.GetToolInput(toolname)
	cobra.CheckErr(err)

	hasErrors, errors := validate.ValidateInputs(toolSpec, toolInput)

	_, err = io.ReadLicenseFile(licenseFile)
	if err != nil {
		warnings = append(warnings, validate.ValidationError{
			Field:    "/src/LICENSE",
			Name:     "LICENSE",
			Message:  "No license file found. We recommend adding one.",
			Type:     validate.ErrorType("WARNING"),
			Expected: "LICENSE",
			Actual:   "None",
		})
	}

	if !hasErrors && len(warnings) == 0 {
		fmt.Println("OK")
		return
	} else if hasErrors {
		fmt.Println("FAIL")
	} else {
		fmt.Println("WARN")
	}

	if len(warnings) > 0 {
		for _, warning := range warnings {
			if verbose {
				fmt.Println("---")
			}
			fmt.Println(io.WriteValidationError(warning, verbose))
		}
	}

	if hasErrors {
		for _, err := range errors {
			if verbose {
				fmt.Println("---")
			}
			fmt.Println(io.WriteValidationError(*err, verbose))
		}
	}

	if verbose && (hasErrors || len(warnings) > 0) {
		fmt.Println("--------------------------------")
		fmt.Printf("ERRORS: %d     WARNINGS: %d\n", len(errors), len(warnings))
		fmt.Println("--------------------------------")
	}
}

func init() {
	verifyCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	rootCmd.AddCommand(verifyCmd)
}
