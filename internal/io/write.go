package io

import (
	"encoding/json"
	"fmt"

	toolspec "github.com/hydrocode-de/tool-spec-go"
)

func InputFileToJSON(input toolspec.InputFile) (string, error) {
	jsonBytes, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal input file to JSON: %w", err)
	}

	return string(jsonBytes), nil
}

func MergeInputFiles(existing, update toolspec.InputFile) toolspec.InputFile {
	merged := existing

	if len(merged) == 0 {
		return update
	}

	for toolname, toolValues := range update {
		toolInput, ok := merged[toolname]
		if !ok {
			toolInput = toolspec.ToolInput{
				Parameters: make(map[string]interface{}),
				Datasets:   make(map[string]string),
			}
		}

		if toolInput.Parameters == nil {
			toolInput.Parameters = make(map[string]interface{})
		}
		if toolInput.Datasets == nil {
			toolInput.Datasets = make(map[string]string)
		}
		for paramName, paramValue := range toolValues.Parameters {
			toolInput.Parameters[paramName] = paramValue
		}
		for dataName, dataValue := range toolValues.Datasets {
			toolInput.Datasets[dataName] = dataValue
		}

		merged[toolname] = toolInput
	}

	return merged
}
