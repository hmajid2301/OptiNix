package options

import (
	"encoding/json"
	"errors"
)

// TODO: see if there is a way we can use this again
type Declaration struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Declarations []Declaration

// INFO: Due to difference in structure between the options files between the different sources.
// A declaration can either be a list of strings (NixOS) or an object with a name and URL for home-manager and darwin.
// This custom unmarshal function handles that for us and allow us to keep the actual parse function a lot cleaner.
func (d *Declarations) UnmarshalJSON(b []byte) error {
	var single Declaration
	if err := json.Unmarshal(b, &single); err == nil {
		*d = Declarations{single}
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
	Loc          []string     `json:"loc"`
	ReadOnly     bool         `json:"readOnly"`
	Type         string       `json:"type"`
}

// TODO: deal with code blocks in examples/descriptions
// i.e. `accounts.calendar.accounts.<name>.remote.passwordCommand`,
// from HM Options
type Option struct {
	Name        string
	Description string
	Type        string
	Default     string
	Example     string
	Sources     []string
}

func ParseOptions(jsonData []byte) ([]Option, error) {
	options := []Option{}
	var jsonOpts map[string]OptionFile
	err := json.Unmarshal(jsonData, &jsonOpts)
	if err != nil {
		return nil, err
	}

	for key, option := range jsonOpts {
		opt := Option{
			Name:        key,
			Description: option.Description,
			Type:        option.Type,
			Example:     option.Example.Text,
			Default:     option.Default.Text,
		}

		for _, declaration := range option.Declarations {
			opt.Sources = append(opt.Sources, declaration.URL)
		}

		options = append(options, opt)
	}
	return options, err
}
