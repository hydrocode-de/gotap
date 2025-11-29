package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hydrocode-de/gotap/internal/config"
	"github.com/hydrocode-de/gotap/internal/input"
	"github.com/hydrocode-de/gotap/internal/io"
	"github.com/hydrocode-de/gotap/internal/validation"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute this tool",
	Long: `Validates input data and starts the tool entrypoint.

The run command optionally creates the input.json, if the parameters are 
provided as command line arguments. Next, the inputs are validated and 
finally the tool is executed.`,
	DisableFlagParsing: true,
	Run:                execute,
}

func execute(cmd *cobra.Command, args []string) {
	var failOnWarnings bool
	for _, arg := range args {
		if arg == "--fail-on-warnings" {
			failOnWarnings = true
		}
	}

	dry, err := PrepareInputs(cmd, args)
	cobra.CheckErr(err)

	result, err := validation.LoadAndValidateSpec(args)
	cobra.CheckErr(err)

	if result.WarningCount() > 0 && failOnWarnings || result.ErrorCount() > 0 {
		fmt.Println("FAIL")

		for _, warning := range result.Warnings {
			fmt.Println(io.WriteValidationError(warning, true))
		}
		for _, err := range result.Errors {
			fmt.Println(io.WriteValidationError(err, true))
		}
		os.Exit(result.ErrorCount())
	}

	command, err := input.ResolveCommand(result.ToolSpec)
	cobra.CheckErr(err)

	if dry {
		fmt.Println(command.Command)
		return
	}

	// execute the command finally. This can later be replaced by
	// by logging, tracing, etc.
	outputFolder := config.GetViper().GetString("output_folder")

	cmdResult, err := input.ExecuteCommand(command)
	cobra.CheckErr(err)

	if cmdResult.Stderr != nil {
		os.WriteFile(filepath.Join(outputFolder, "STDERR"), cmdResult.Stderr, 0644)
	}
	if cmdResult.Stdout != nil {
		os.WriteFile(filepath.Join(outputFolder, "STDOUT"), cmdResult.Stdout, 0644)
	}
	jsonResult, err := json.MarshalIndent(cmdResult, "", "  ")
	if err == nil {
		os.WriteFile(filepath.Join(outputFolder, "_metadata.json"), jsonResult, 0644)
	}
}

func init() {
	runCmd.Flags().Bool("dry", false, "Dry run the tool, returning the new inputs.json, instead of executing the tool.")
	runCmd.Flags().Bool("update-inputs", false, "Update the inputs.json if arguments are provided and the file already exists.")

	runCmd.Flags().Bool("fail-on-warnings", false, "Fail the tool if there are warnings.")
	rootCmd.AddCommand(runCmd)
}
