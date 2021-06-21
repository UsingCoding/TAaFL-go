package app

import (
	"compiler/pkg/common/grammary"
	"compiler/pkg/slr/common/inlinedgrammary"
)

const (
	newAxiomSymbol = "|F|"
)

func NewFilledGrammarFilter() InlinedGrammarFilter {
	return &filledGrammarFilter{}
}

type filledGrammarFilter struct {
}

func (filter *filledGrammarFilter) Filter(grammar inlinedgrammary.Grammar) (inlinedgrammary.Grammar, error) {
	rules := grammar.FindRulesThatContains(grammar.Axiom())
	if len(rules) == 0 {
		return grammar, nil
	}

	axiomSymbol := grammary.NewSymbol(newAxiomSymbol)

	newRules := append(
		[]inlinedgrammary.Rule{inlinedgrammary.NewRule(axiomSymbol, []grammary.Symbol{grammar.Axiom()})},
		grammar.Rules()...,
	)

	return inlinedgrammary.New(
		axiomSymbol,
		newRules...,
	), nil
}
