package config

import (
	"fmt"
	"os"

	toolspec "github.com/hydrocode-de/tool-spec-go"
)

func ResolveToolname(args []string, inputs toolspec.InputFile) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}

	if toolname := os.Getenv("RUN_TOOL"); toolname != "" {
		return toolname, nil
	}

	if len(inputs) == 1 {
		for toolname := range inputs {
			return toolname, nil
		}
	}

	return "", fmt.Errorf("the toolname could not be resolved. Pass it as an argument or set the RUN_TOOL environment variable")
}
