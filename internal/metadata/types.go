package metadata

import toolspec "github.com/hydrocode-de/tool-spec-go"

type Converter interface {
	Ingest(spec toolspec.ToolSpec)
	Validate() bool
	Serialize(format string) ([]byte, error)
}
