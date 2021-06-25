package app

import (
	"compiler/pkg/common/grammary"
	"compiler/pkg/slr/common/inlinedgrammary"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type InlinedGrammarParser interface {
	Parse(rawGrammar string) (inlinedgrammary.Grammar, error)
}

func NewInlinedGrammarParser() InlinedGrammarParser {
	return &parser{}
}

type parser struct {
}

func (p *parser) Parse(rawGrammar string) (inlinedgrammary.Grammar, error) {
	lines := strings.Split(rawGrammar, "\n")

	var axiom *grammary.Symbol
	rules := make([]inlinedgrammary.Rule, 0, len(lines))

	for i, line := range lines {
		ruleParts := strings.Split(line, "->")

		if len(ruleParts) != 2 {
			return inlinedgrammary.Grammar{}, errors.New(fmt.Sprintf("unknown grammar format at line %d", i))
		}

		leftSideSymbol := grammary.NewSymbol(strings.Trim(ruleParts[0], " "))
		if axiom == nil {
			axiom = &leftSideSymbol
		}

		symbols := strings.Split(strings.Trim(ruleParts[1], " "), " ")

		ruleSymbols := make([]grammary.Symbol, 0, len(symbols))

		for _, symbol := range symbols {
			ruleSymbols = append(ruleSymbols, grammary.NewSymbol(symbol))
		}

		rules = append(rules, inlinedgrammary.NewRule(leftSideSymbol, ruleSymbols))
	}

	if axiom == nil {
		return inlinedgrammary.Grammar{}, errors.New("no axiom found, empty grammar")
	}

	return inlinedgrammary.New(*axiom, rules...), nil
}
