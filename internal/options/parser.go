package options

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type Option struct {
	Name        string
	Description string
	Type        string
	Default     string
	Source      string
}

var ignore = map[string]struct{}{"_module.args": {}}

func Parse(node *html.Node) []Option {
	options := []Option{}
	doc := goquery.NewDocumentFromNode(node)
	doc.Find("dt").Each(func(_ int, s *goquery.Selection) {
		option := getOption(s)
		if option == (Option{}) {
			return
		}

		options = append(options, option)
	})
	return options
}

func getOption(tableHeader *goquery.Selection) Option {
	option := Option{}
	optionName := strings.TrimSpace(tableHeader.Text())
	if _, ok := ignore[optionName]; ok {
		return option
	}
	option.Name = optionName

	tableElements := tableHeader.NextFiltered("dd")
	optionFields := tableElements.Find("p")
	optionFields.Each(func(_ int, d *goquery.Selection) {
		option.getOptionField(d)
	})
	return option
}

func (o *Option) getOptionField(optionParagraph *goquery.Selection) {
	optionData := getOptionData(optionParagraph)

	// INFO: field title and data are one string, need to split by colon (:) to figure out the title.
	// e.g. ""
	optionFields := strings.Split(optionData, ":")

	switch optionTitle := strings.TrimSpace(optionFields[0]); {
	case optionTitle == "Type":
		optionContents := strings.TrimSpace(optionFields[1])
		o.Type = optionContents
	case optionTitle == "Declared by":
		sourceAnchorLink := optionParagraph.SiblingsFiltered("table").First().Find("a.filename")
		o.Source = sourceAnchorLink.AttrOr("href", "")
	case optionTitle == "Default":
		optionContents := strings.TrimSpace(optionFields[1])
		o.Default = optionContents
	default:
		optionContents := strings.TrimSpace(optionData)
		o.Description = optionContents
	}
}

func getOptionData(optionParagraph *goquery.Selection) string {
	var optionData string
	optionParagraph.Contents().Each(func(_ int, optionChild *goquery.Selection) {
		if optionChild.Is("a") {
			link := optionChild.AttrOr("href", "")
			optionData += "[" + optionChild.Text() + "](" + link + ")"
			return
		}

		optionData += optionChild.Text()
	})
	return optionData
}
