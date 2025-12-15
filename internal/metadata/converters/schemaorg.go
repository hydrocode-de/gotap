package converters

import (
	"encoding/json"
	"fmt"

	toolspec "github.com/hydrocode-de/tool-spec-go"
)

type SchemaOrg struct {
	Context             string          `json:"@context"`
	Type                string          `json:"@type"`
	Name                string          `json:"name,omitempty"`
	Description         string          `json:"description,omitempty"`
	SoftwareVersion     string          `json:"softwareVersion,omitempty"`
	Author              []Person        `json:"author,omitempty"`
	License             string          `json:"license,omitempty"`
	CodeRepository      string          `json:"codeRepository,omitempty"`
	Identifier          []PropertyValue `json:"identifier,omitempty"`
	Keywords            []string        `json:"keywords,omitempty"`
	ApplicationCategory string          `json:"applicationCategory,omitempty"`
	OperationgSystem    []string        `json:"operatingSystem,omitempty"`
	ProgrammingLanguage []string        `json:"programmingLanguage,omitempty"`
	AdditionalProperty  []PropertyValue `json:"additionalProperty,omitempty"`
}

type Person struct {
	Type        string          `json:"@type"`
	Name        string          `json:"name,omitempty"`
	Email       string          `json:"email,omitempty"`
	URL         string          `json:"url,omitempty"`
	SameAs      string          `json:"sameAs,omitempty"`
	Affiliation string          `json:"affiliation,omitempty"`
	Identifier  []PropertyValue `json:"identifier,omitempty"`
	JobTitle    string          `json:"jobTitle,omitempty"`
}

type Organization struct {
	Type string `json:"@type"`
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type PropertyValue struct {
	Type       string `json:"@type"`
	PropertyID string `json:"propertyID,omitempty"`
	Value      string `json:"value,omitempty"`
}

type SchemaOrgConverter struct {
	SchemaOrg
	errs []error
}

func (s *SchemaOrgConverter) Ingest(spec toolspec.ToolSpec) {
	if spec.Citation.Title == "" {
		s.errs = append(s.errs, fmt.Errorf("critical. no citation found"))
	}
	s.SchemaOrg = SchemaOrg{
		Context:     "https://schema.org",
		Type:        "SoftwareApplication",
		Name:        spec.Title,
		Description: spec.Description,
	}

	if len(spec.Citation.Authors) > 0 {
		for _, author := range spec.Citation.Authors {
			var person Person
			if author.IsPerson {
				person = Person{
					Type:   "Person",
					Name:   author.Person.Family + ", " + author.Person.GivenNames,
					Email:  author.Person.Email,
					URL:    author.Person.Website.URL.String(),
					SameAs: string(author.Person.Orcid),
				}
				if len(author.Person.Affiliation) > 0 {
					person.Affiliation = author.Person.Affiliation
				}
			} else {
				person = Person{
					Type: "Organization",
					Name: author.Entity.Name,
					URL:  author.Entity.Website.URL.String(),
				}
			}
			s.SchemaOrg.Author = append(s.SchemaOrg.Author, person)
		}
	}

	for name, param := range spec.Parameters {
		desc := param.Description
		if desc == "" {
			desc = fmt.Sprintf("%s (%s)", name, param.ToolType)
		}
		prop := PropertyValue{
			Type:       "PropertyValue",
			PropertyID: name,
			Value:      desc,
		}
		s.SchemaOrg.AdditionalProperty = append(s.SchemaOrg.AdditionalProperty, prop)
	}
}

func (s *SchemaOrgConverter) Validate() bool {
	return s.errs == nil
}

func (s *SchemaOrgConverter) Serialize(format string) ([]byte, error) {
	return json.MarshalIndent(s.SchemaOrg, "", "  ")
}
