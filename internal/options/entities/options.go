package entities

type Option struct {
	Name        string
	Description string
	Type        string
	Default     string
	Example     string
	OptionFrom  string
	Sources     []string
}

type Sources struct {
	NixOS       string
	HomeManager string
	Darwin      string
}
