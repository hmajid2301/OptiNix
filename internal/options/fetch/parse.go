package fetch

import (
	"encoding/json"
	"errors"

	"gitlab.com/hmajid2301/optinix/internal/options/entities"
)

type Declaration struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Declarations []Declaration

// INFO: Due to difference in structure between the options files between the different sources.
// A declaration can either be a list of strings (NixOS) or an object with a name and URL for home-manager and darwin.
// This custom unmarshal function handles that for us and allow us to keep the actual parse function a lot cleaner.
func (d *Declarations) UnmarshalJSON(b []byte) error {
	var declarations []Declaration
	if err := json.Unmarshal(b, &declarations); err == nil {
		*d = declarations
		return nil
	}

	var slice []string
	if err := json.Unmarshal(b, &slice); err == nil {
		for _, s := range slice {
			*d = append(*d, Declaration{URL: s})
		}
		return nil
	}

	return errors.New("failed to unmarshal declarations")
}

type Default struct {
	Type string `json:"_type"`
	Text string `json:"text"`
}

type Example struct {
	Type string `json:"_type"`
	Text string `json:"text"`
}

type OptionFile struct {
	Declarations Declarations `json:"declarations"`
	Default      Default      `json:"default"`
	Description  string       `json:"description"`
	Example      Example      `json:"example"`
	Type         string       `json:"type"`
	Loc          []string     `json:"loc"`
	ReadOnly     bool         `json:"readOnly"`
}

func ParseOptions(jsonData []byte, optionFrom string) ([]entities.Option, error) {
	options := []entities.Option{}
	var jsonOpts map[string]OptionFile
	err := json.Unmarshal(jsonData, &jsonOpts)
	if err != nil {
		return nil, err
	}

	for key, option := range jsonOpts {
		opt := entities.Option{
			Name:        key,
			Description: option.Description,
			Type:        option.Type,
			Example:     option.Example.Text,
			Default:     option.Default.Text,
			OptionFrom:  optionFrom,
		}

		for _, declaration := range option.Declarations {
			opt.Sources = append(opt.Sources, declaration.URL)
		}

		options = append(options, opt)
	}
	return options, err
}
