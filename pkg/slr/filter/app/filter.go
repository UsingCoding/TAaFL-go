package app

import (
	"compiler/pkg/slr/common/inlinedgrammary"
)

type InlinedGrammarFilter interface {
	Filter(grammar inlinedgrammary.Grammar) (inlinedgrammary.Grammar, error)
}
