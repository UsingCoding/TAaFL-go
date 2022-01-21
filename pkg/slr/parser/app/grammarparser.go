package app

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"compiler/pkg/common/grammary"
	"compiler/pkg/slr/common"
	"compiler/pkg/slr/common/inlinedgrammary"
)

var (
	astRuleRegexp = regexp.MustCompile(`<<.*>>`)
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

		var astRulePtr *common.ASTRule

		if rawASTRule := symbols[len(symbols)-1]; astRuleRegexp.MatchString(rawASTRule) {
			rawASTRule = strings.Trim(rawASTRule, "<<>>")
			astRulePtr = (*common.ASTRule)(&rawASTRule)

			symbols = symbols[:len(symbols)-1]
		}

		rule := inlinedgrammary.NewRule(leftSideSymbol, ruleSymbols)
		if astRulePtr != nil {
			rule.SetASTRule(*astRulePtr)
		}
		rules = append(rules, rule)
	}

	if axiom == nil {
		return inlinedgrammary.Grammar{}, errors.New("no axiom found, empty grammar")
	}

	return inlinedgrammary.New(*axiom, rules...), nil
}
