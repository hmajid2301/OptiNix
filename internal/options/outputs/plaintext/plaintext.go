package plaintext

import (
	"fmt"

	"gitlab.com/hmajid2301/optinix/internal/options/entities"
)

func Output(options []entities.Option) {
	for idx, opt := range options {
		fmt.Printf("\n\nOption: %d\n", idx)
		fmt.Printf("Name: %s\n", opt.Name)
		fmt.Printf("Type: %s\n", opt.Type)
		fmt.Printf("Default: %s\n", opt.Default)
		fmt.Printf("Example: %s\n", opt.Example)
		fmt.Printf("From: %s\n", opt.OptionFrom)
		fmt.Printf("Sources: %s\n", opt.Sources)
	}
}
