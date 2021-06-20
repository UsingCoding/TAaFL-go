package inlinedgrammary

import "compiler/pkg/common/grammary"

type Grammar struct {
	axiom grammary.Symbol
	rules []Rule
}

func (g *Grammar) Axiom() grammary.Symbol {
	return g.axiom
}

func (g *Grammar) Rules() []Rule {
	return g.rules
}

func (g *Grammar) FindByLeftSideSymbol(symbol grammary.Symbol) map[uint]Rule {
	result := map[uint]Rule{}
	for ruleNumber, rule := range g.rules {
		if rule.leftSideSymbol == symbol {
			result[uint(ruleNumber)] = rule
		}
	}
	return result
}

func (g *Grammar) FindRulesThatContains(needle grammary.Symbol) map[uint]Rule {
	result := map[uint]Rule{}
	for ruleNumber, rule := range g.rules {
		for _, symbol := range rule.RuleSymbols() {
			if needle == symbol {
				result[uint(ruleNumber)] = rule
			}
		}
	}
	return result
}

type Rule struct {
	leftSideSymbol grammary.Symbol
	ruleSymbols    []grammary.Symbol
}

func (r *Rule) LeftSideSymbol() grammary.Symbol {
	return r.leftSideSymbol
}

func (r *Rule) RuleSymbols() []grammary.Symbol {
	return r.ruleSymbols
}

func (r *Rule) Find(symbol grammary.Symbol) *int {
	for i, ruleSymbol := range r.ruleSymbols {
		if ruleSymbol == symbol {
			return &i
		}
	}
	return nil
}

func New(axiom grammary.Symbol, rules ...Rule) Grammar {
	return Grammar{
		axiom: axiom,
		rules: rules,
	}
}

// NewRule TODO: remove after implementing table generator
func NewRule(leftSideSymbol grammary.Symbol, symbols []grammary.Symbol) Rule {
	return Rule{
		leftSideSymbol: leftSideSymbol,
		ruleSymbols:    symbols,
	}
}

func NewFromGrammar(grammar grammary.Grammar) Grammar {
	grammarAxiom := *grammar.Axiom

	newGrammar := Grammar{
		axiom: grammarAxiom,
	}

	for _, roll := range grammar.Impl[grammarAxiom] {
		var dstRoll []grammary.Symbol
		copy(dstRoll, roll)
		newGrammar.rules = append(newGrammar.rules, Rule{
			leftSideSymbol: grammarAxiom,
			ruleSymbols:    dstRoll,
		})
	}

	for leftSideSymbol, rolls := range grammar.Impl {
		if leftSideSymbol == grammarAxiom {
			continue
		}

		for _, roll := range rolls {
			var dstRoll []grammary.Symbol
			copy(dstRoll, roll)
			newGrammar.rules = append(newGrammar.rules, Rule{
				leftSideSymbol: leftSideSymbol,
				ruleSymbols:    dstRoll,
			})
		}
	}

	return newGrammar
}
