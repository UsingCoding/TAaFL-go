package grammary

import (
	"github.com/pkg/errors"
	"regexp"
	"strings"
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

func IsNonTerminalSymbol(value string) bool {
	matched, _ := regexp.MatchString(`<[A-Z]>`, value)
	return matched
}

func Parse(rawData string) (Grammar, error) {
	grammar := NewGrammar()

	for _, line := range strings.Split(rawData, "\n") {
		// skip comments strings
		if strings.HasPrefix(line, "#") {
			continue
		}
		var buffer string

		var alternatives [][]Symbol
		var currentAlternativeRoll []Symbol
		var leftSideSymbol *Symbol
		leftSide := true

		for _, i := range line {
			ch := string(i)
			if ch != " " {
				buffer += ch
				continue
			}

			if buffer == RuleSidesSeparator {
				if leftSideSymbol == nil {
					return Grammar{}, errors.New("no leftSideSymbol found")
				}

				if !leftSide {
					return Grammar{}, errors.New("you should escape -> with \\")
				}

				buffer = ""
				leftSide = false
				continue
			}

			if IsNonTerminalSymbol(buffer) && leftSideSymbol == nil {
				newSymbol := NewSymbol(buffer)
				leftSideSymbol = &newSymbol
				buffer = ""
				// check if we already this rule and this another alternative
				alts, exists := grammar.Impl[newSymbol]
				if exists {
					alternatives = alts
				}
				continue
			}

			if buffer == "|" {
				if len(currentAlternativeRoll) == 0 {
					return Grammar{}, errors.New("no symbols before alternative")
				}

				alternatives = append(alternatives, currentAlternativeRoll)
				currentAlternativeRoll = nil
				buffer = ""
				continue
			}

			currentAlternativeRoll = append(currentAlternativeRoll, NewSymbol(buffer))
			buffer = ""
		}

		if buffer != "" {
			currentAlternativeRoll = append(currentAlternativeRoll, NewSymbol(buffer))
		}

		if leftSideSymbol == nil {
			return Grammar{}, errors.New("rule on left side not found")
		}

		if len(currentAlternativeRoll) != 0 {
			alternatives = append(alternatives, currentAlternativeRoll)
		}

		grammar.AddRule(*leftSideSymbol, alternatives)

		if grammar.Axiom == nil {
			grammar.Axiom = leftSideSymbol
		}
	}

	return grammar, nil
}
