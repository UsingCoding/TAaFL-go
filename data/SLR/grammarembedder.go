package SLR

import _ "embed"

//go:embed grammar
var grammar string

func EmbeddedGrammar() string {
	return grammar
}
