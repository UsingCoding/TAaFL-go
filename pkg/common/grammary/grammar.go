package grammary

import (
	"regexp"
)

const (
	EmptySymbol   = "e"
	EndOfSequence = "_|_"

	RuleSidesSeparator    = "->"
	RuleSequenceSeparator = "/"
)

func NewSymbol(ch string) Symbol {
	return Symbol{ch: ch}
}

type Symbol struct {
	ch string
}

func (s Symbol) NonTerminal() bool {
	return IsNonTerminalSymbol(s.ch)
}

func (s Symbol) String() string {
	return s.ch
}

func NewGrammar() Grammar {
	return Grammar{
		Impl: make(map[Symbol][][]Symbol),
	}
}

type Grammar struct {
	Axiom *Symbol
	Impl  map[Symbol][][]Symbol
}

func (g *Grammar) AddRule(ruleLeftSide Symbol, alternatives [][]Symbol) {
	g.Impl[ruleLeftSide] = alternatives
}

func (g *Grammar) Copy() Grammar {
	dest := Grammar{
		Axiom: *(&g.Axiom),
	}

	for symbol, rolls := range g.Impl {
		var destRolls [][]Symbol

		for _, roll := range rolls {
			var destRoll []Symbol
			copy(destRoll, roll)
			destRolls = append(destRolls, destRoll)
		}

		dest.Impl[symbol] = destRolls
	}

	return dest
}

func IsNonTerminalSymbol(value string) bool {
	matched, _ := regexp.MatchString(`<[A-Za-z]*>`, value)
	return matched
}
