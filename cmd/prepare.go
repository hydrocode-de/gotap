/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/hydrocode-de/gotap/internal/config"
	"github.com/hydrocode-de/gotap/internal/input"
	"github.com/hydrocode-de/gotap/internal/io"
	"github.com/hydrocode-de/gotap/internal/validation"
	toolspec "github.com/hydrocode-de/tool-spec-go"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// prepareCmd represents the prepare command
var prepareCmd = &cobra.Command{
	Use:   "prepare",
	Short: "Create tool execution input file",
	Long: `Creates the inputs.json parameterization file.

The tool execution relies on all parameters and datasets being statically
defined in the inputs.json file. Using prepare --help you can learn about
all available parameters and datasets settings and use the cli to create
or update the inputs.json, before running it.`,
	Run:                prepare,
	DisableFlagParsing: true,
}

func prepare(cmd *cobra.Command, args []string) {
	_, err := PrepareInputs(cmd, args)
	cobra.CheckErr(err)
	os.Exit(0)
}

func PrepareInputs(cmd *cobra.Command, args []string) (bool, error) {
	var remainingArgs []string
	specFilePath := ""
	inputFilePath := ""
	updateInputs := false
	dry := false

	for i := 0; i < len(args); i++ {
		if args[i] == "--spec-file" && i+1 < len(args) {
			specFilePath = args[i+1]
			i++ // Skip the value
			continue
		}
		if args[i] == "--input-file" && i+1 < len(args) {
			inputFilePath = args[i+1]
			i++ // Skip the value
			continue
		}
		if args[i] == "--update-inputs" {
			updateInputs = true
			continue
		}
		if args[i] == "--dry" {
			dry = true
			continue
		}
		remainingArgs = append(remainingArgs, args[i])
	}

	if specFilePath != "" {
		config.GetViper().Set("spec_file", specFilePath)
	}
	if inputFilePath != "" {
		config.GetViper().Set("input_file", inputFilePath)
	}

	toolSpec, err := validation.LoadSpec(remainingArgs)
	if err != nil {
		return dry, err
	}

	// remove the toolname from the args
	if len(remainingArgs) > 0 && remainingArgs[0] == toolSpec.Name {
		remainingArgs = remainingArgs[1:]
	}

	for _, arg := range remainingArgs {
		if arg == "-h" || arg == "--help" {
			input.RegisterFlags(toolSpec, cmd.Flags())
			cmd.Help()
			os.Exit(0)
		}
	}

	flagSet := pflag.NewFlagSet("run", pflag.ContinueOnError)
	input.RegisterFlags(toolSpec, flagSet)

	err = flagSet.Parse(remainingArgs)
	if err != nil {
		if strings.Contains(err.Error(), "help requested") {
			return dry, nil
		}
		return dry, err

	}

	inputFile, err := input.CollectInputs(toolSpec, flagSet)
	if err != nil {
		return dry, err
	}

	numDynParams := len(inputFile[toolSpec.Name].Parameters)
	numDynData := len(inputFile[toolSpec.Name].Datasets)
	hasDynamicFlags := numDynParams+numDynData > 0

	outputPath := config.GetViper().GetString("input_file")
	_, err = os.Stat(outputPath)
	fileExists := err == nil

	var toolInput toolspec.InputFile
	if fileExists {
		inputValues, err := io.ReadInputFile(outputPath)
		if err != nil {
			return dry, err
		}
		toolInput = io.MergeInputFiles(inputValues, inputFile)
	} else {
		toolInput = inputFile
	}

	jsonInput, err := io.InputFileToJSON(toolInput)
	if err != nil {
		return dry, err
	}

	if dry {
		fmt.Printf("%s\n", jsonInput)
		return dry, nil
	}

	if hasDynamicFlags && fileExists && !updateInputs {
		return dry, fmt.Errorf("inputs.json already exists. Use --update-inputs to update the file")
	}
	if hasDynamicFlags {
		err = os.WriteFile(outputPath, []byte(jsonInput), 0644)
		if err != nil {
			return dry, err
		}
		return dry, nil
	}

	return dry, nil
}

func init() {
	prepareCmd.Flags().Bool("dry", false, "Dry run the tool, returning the new inputs.json, instead of executing the tool.")
	prepareCmd.Flags().Bool("update-inputs", false, "Update the inputs.json if arguments are provided and the file already exists.")

	rootCmd.AddCommand(prepareCmd)

}
