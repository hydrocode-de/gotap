package input

import (
	toolspec "github.com/hydrocode-de/tool-spec-go"
	"github.com/spf13/pflag"
)

func RegisterFlags(spec toolspec.ToolSpec, flagSet *pflag.FlagSet) {
	// Build the flag set - no defaults, params must be set explicitly
	for paramName, param := range spec.Parameters {
		switch param.ToolType {
		case "string", "enum", "date", "datetime", "time":
			if param.IsArray {
				flagSet.StringSlice(paramName, nil, param.Description)
			} else {
				flagSet.String(paramName, "", param.Description)
			}
		case "integer":
			if param.IsArray {
				flagSet.IntSlice(paramName, nil, param.Description)
			} else {
				flagSet.Int(paramName, 0, param.Description)
			}
		case "float":
			if param.IsArray {
				flagSet.Float64Slice(paramName, nil, param.Description)
			} else {
				flagSet.Float64(paramName, 0.0, param.Description)
			}
		case "boolean":
			if param.IsArray {
				flagSet.BoolSlice(paramName, nil, param.Description)
			} else {
				flagSet.Bool(paramName, false, param.Description)
			}
		}
	}
	for dataName, data := range spec.Data {
		flagSet.String(dataName, "", data.Description)
	}
}

func CollectInputs(spec toolspec.ToolSpec, flagSet *pflag.FlagSet) (toolspec.InputFile, error) {
	inputParameters := make(map[string]interface{})
	inputData := make(map[string]string)

	flagSet.Visit(func(f *pflag.Flag) {
		// Only process flags that were actually set
		param := spec.Parameters[f.Name]

		switch param.ToolType {
		case "string", "enum", "date", "datetime", "time":
			if param.IsArray {
				val, _ := flagSet.GetStringSlice(f.Name)
				inputParameters[f.Name] = val
			} else {
				val, _ := flagSet.GetString(f.Name)
				inputParameters[f.Name] = val
			}
		case "integer":
			if param.IsArray {
				val, _ := flagSet.GetIntSlice(f.Name)
				inputParameters[f.Name] = val
			} else {
				val, _ := flagSet.GetInt(f.Name)
				inputParameters[f.Name] = val
			}
		case "float":
			if param.IsArray {
				val, _ := flagSet.GetFloat64Slice(f.Name)
				inputParameters[f.Name] = val
			} else {
				val, _ := flagSet.GetFloat64(f.Name)
				inputParameters[f.Name] = val
			}
		case "boolean":
			if param.IsArray {
				val, _ := flagSet.GetBoolSlice(f.Name)
				inputParameters[f.Name] = val
			} else {
				val, _ := flagSet.GetBool(f.Name)
				inputParameters[f.Name] = val
			}
		}

		// handle data
		if _, ok := spec.Data[f.Name]; ok {
			val, _ := flagSet.GetString(f.Name)
			inputData[f.Name] = val
		}
	})

	inputFile := toolspec.InputFile{}
	inputFile[spec.Name] = toolspec.ToolInput{
		Parameters: inputParameters,
		Datasets:   inputData,
	}

	return inputFile, nil
}
