package plaintext

import (
	"fmt"
	"strings"

	"gitlab.com/hmajid2301/optinix/internal/options/entities"
)

func Output(options []entities.Option) {
	for idx, opt := range options {
		separator := strings.Repeat("─", 80)

		fmt.Printf("\n%s\n", separator)
		fmt.Printf("Option %d: %s\n", idx+1, opt.Name)
		fmt.Printf("%s\n\n", separator)

		if opt.Description != "" {
			wrapped := wrapText(opt.Description, 76)
			fmt.Printf("Description:\n  %s\n\n", wrapped)
		}

		fmt.Printf("Type:    %s\n", opt.Type)

		if opt.Default != "" {
			fmt.Printf("Default: %s\n", opt.Default)
		}

		if opt.Example != "" {
			fmt.Printf("Example: %s\n", opt.Example)
		}

		fmt.Printf("Source:  %s\n", opt.OptionFrom)

		if len(opt.Sources) > 0 {
			fmt.Printf("\nDefined in:\n")
			for _, source := range opt.Sources {
				fmt.Printf("  - %s\n", source)
			}
		}
	}

	if len(options) > 0 {
		fmt.Printf("\n%s\n", strings.Repeat("─", 80))
		fmt.Printf("Total: %d option(s)\n", len(options))
	}
}

func wrapText(text string, width int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		if currentLine.Len() == 0 {
			currentLine.WriteString(word)
		} else if currentLine.Len()+1+len(word) <= width {
			currentLine.WriteString(" ")
			currentLine.WriteString(word)
		} else {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
			currentLine.WriteString(word)
		}
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return strings.Join(lines, "\n  ")
}
