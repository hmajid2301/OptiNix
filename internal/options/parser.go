package options

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// TODO: deal with code blocks in examples/descriptions
// i.e. `accounts.calendar.accounts.<name>.remote.passwordCommand`,
// from https://nix-community.github.io/home-manager/options.html.
type Option struct {
	Name        string
	Description string
	Type        string
	Default     string
	Example     string
	Sources     []string
}

var ignore = map[string]struct{}{"_module.args": {}}

func Parse(node *html.Node) []Option {
	options := []Option{}
	doc := goquery.NewDocumentFromNode(node)
	doc.Find("dt").Each(
		func(_ int, dt *goquery.Selection) {
			option := getOption(dt)
			if option.Name == "" {
				return
			}
			options = append(options, option)
		},
	)

	return options
}

func getOption(dt *goquery.Selection) Option {
	option := Option{}
	optionName := strings.TrimSpace(dt.Text())
	if _, ok := ignore[optionName]; ok {
		return option
	}
	option.Name = optionName

	dt.NextFiltered("dd").Find("p").Each(
		func(_ int, p *goquery.Selection) {
			option.setOptionField(p)
		},
	)
	return option
}

func (o *Option) setOptionField(p *goquery.Selection) {
	optionData := getContent(p)

	// INFO: field title and data are one string, need to split by colon (:) to figure out the title.
	// e.g. ""
	optionFields := strings.Split(optionData, ":")

	switch optionTitle := strings.TrimSpace(optionFields[0]); {
	case optionTitle == "Type":
		optionContents := strings.TrimSpace(optionFields[1])
		o.Type = optionContents
	case optionTitle == "Example":
		optionContents := strings.TrimSpace(optionFields[1])
		o.Example = optionContents
	case optionTitle == "Declared by":
		p.SiblingsFiltered("table").First().Find("a.filename").Each(func(_ int, anchors *goquery.Selection) {
			o.Sources = append(o.Sources, anchors.AttrOr("href", ""))
		})
	case optionTitle == "Default":
		optionContents := strings.TrimSpace(optionFields[1])
		o.Default = optionContents
	default:
		optionContents := strings.TrimSpace(optionData)
		o.Description = optionContents
	}
}

func getContent(p *goquery.Selection) string {
	var optionContent string
	p.Contents().Each(func(_ int, optionData *goquery.Selection) {
		if optionData.Is("a") {
			link := optionData.AttrOr("href", "")
			optionContent += "[" + optionData.Text() + "](" + link + ")"
			return
		}

		optionContent += optionData.Text()
	})
	return optionContent
}
