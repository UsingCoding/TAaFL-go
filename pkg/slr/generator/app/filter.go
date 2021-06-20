package app

import (
	"compiler/pkg/common/grammary"
	"compiler/pkg/slr/common/inlinedgrammary"
	"fmt"
	"github.com/pkg/errors"
)

type InlinedGrammarFilter interface {
	Filter(grammar inlinedgrammary.Grammar) (inlinedgrammary.Grammar, error)
}

func NewEmptySymbolFilter() InlinedGrammarFilter {
	return &emptySymbolFilter{}
}

type emptySymbolFilter struct {
}

func (filter *emptySymbolFilter) Filter(grammar inlinedgrammary.Grammar) (newGrammar inlinedgrammary.Grammar, err error) {
	emptySymbol := grammary.NewSymbol(grammary.EmptySymbol)

	newGrammar = grammar

	for _, rule := range newGrammar.Rules() {
		ruleSymbols := rule.RuleSymbols()
		if len(ruleSymbols) == 1 && ruleSymbols[0] == emptySymbol {
			ruleLeftSymbol := rule.LeftSideSymbol()
			filterableRulesMap := newGrammar.FindRulesThatContains(ruleLeftSymbol)

			if len(filterableRulesMap) == 0 {
				continue
			}

			newGrammar, err = filter.populateWithFilteredRules(
				newGrammar,
				filterableRulesMap,
				ruleLeftSymbol,
			)
			if err != nil {
				return inlinedgrammary.Grammar{}, err
			}
		}
	}

	for ruleNumber, rule := range newGrammar.Rules() {
		ruleSymbols := rule.RuleSymbols()
		if len(ruleSymbols) == 1 && ruleSymbols[0] == emptySymbol {
			newGrammar, err = filter.removeRuleFromGrammar(newGrammar, ruleNumber)
			if err != nil {
				return inlinedgrammary.Grammar{}, err
			}
		}
	}

	// temporary commented
	//sort.SliceStable(newGrammar.Rules(), func(i, j int) bool {
	//	a := newGrammar.Rules()[i].LeftSideSymbol()
	//	b := newGrammar.Rules()[j].LeftSideSymbol()
	//	if a.String() < b.String() {
	//		return false
	//	}
	//	return true
	//})

	return newGrammar, nil
}

func (filter *emptySymbolFilter) populateWithFilteredRules(
	grammar inlinedgrammary.Grammar,
	filterableRules map[uint]inlinedgrammary.Rule,
	symbolToFilter grammary.Symbol,
) (inlinedgrammary.Grammar, error) {
	newRules := make([]inlinedgrammary.Rule, 0, len(filterableRules))
	for _, rule := range filterableRules {
		pos := rule.Find(symbolToFilter)
		if pos == nil {
			return inlinedgrammary.Grammar{}, errors.New(fmt.Sprintf("expected symbol not in rule : %s : %s", symbolToFilter, rule.RuleSymbols()))
		}

		dstSymbols := make([]grammary.Symbol, len(rule.RuleSymbols()))
		srcSymbols := rule.RuleSymbols()

		copy(dstSymbols, srcSymbols)

		newRules = append(newRules, inlinedgrammary.NewRule(
			rule.LeftSideSymbol(),
			append(dstSymbols[:*pos], dstSymbols[*pos+1:]...),
		))
	}

	return inlinedgrammary.New(
		grammar.Axiom(),
		append(grammar.Rules(), newRules...)...,
	), nil
}

func (filter *emptySymbolFilter) removeFromRule(rule inlinedgrammary.Rule, needle grammary.Symbol) inlinedgrammary.Rule {
	for numberInRule, symbol := range rule.RuleSymbols() {
		if needle == symbol {
			return inlinedgrammary.NewRule(
				rule.LeftSideSymbol(),
				append(rule.RuleSymbols()[:numberInRule], rule.RuleSymbols()[numberInRule+1:]...),
			)
		}
	}
	return rule
}

func (filter *emptySymbolFilter) removeRuleFromGrammar(grammar inlinedgrammary.Grammar, ruleNumber int) (inlinedgrammary.Grammar, error) {
	if ruleNumber >= len(grammar.Rules()) {
		return inlinedgrammary.Grammar{}, errors.WithStack(errors.New("cant remove not existing rule"))
	}

	return inlinedgrammary.New(
		grammar.Axiom(),
		append(grammar.Rules()[:ruleNumber], grammar.Rules()[ruleNumber+1:]...)...,
	), nil
}
