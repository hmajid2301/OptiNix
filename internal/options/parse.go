package options

import (
	"encoding/json"
)

type Declaration struct {
	Name string `json:"name"`
	URL  string `json:"url"`
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
	Declarations []Declaration `json:"declarations"`
	Default      Default       `json:"default"`
	Description  string        `json:"description"`
	Example      Example       `json:"example"`
	Loc          []string      `json:"loc"`
	ReadOnly     bool          `json:"readOnly"`
	Type         string        `json:"type"`
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
		sources := []string{}
		for _, declaration := range option.Declarations {
			sources = append(sources, declaration.URL)
		}
		opt := Option{
			Name:        key,
			Description: option.Description,
			Type:        option.Type,
			Example:     option.Example.Text,
			Default:     option.Default.Text,
		}
		opt.Sources = sources
		options = append(options, opt)
	}
	return options, err
}
