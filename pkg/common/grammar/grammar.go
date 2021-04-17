package grammar

import (
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

const ruleSidesSeparator = "->"

func NewSymbol(ch string) Symbol {
	return Symbol{ch: ch}
}

type Symbol struct {
	ch string
}

func (s Symbol) Terminal() bool {
	return IsTerminalSymbol(s.ch)
}

func IsTerminalSymbol(value string) bool {
	matched, _ := regexp.MatchString(`<[A-Z]>`, value)
	return matched
}

func NewGrammar() Grammar {
	return Grammar{
		impl: make(map[Symbol][][]Symbol),
	}
}

type Grammar struct {
	Axiom *Symbol
	impl  map[Symbol][][]Symbol
}

func (g *Grammar) AddRule(ruleLeftSide Symbol, alternatives [][]Symbol) {
	g.impl[ruleLeftSide] = alternatives
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

			if buffer == ruleSidesSeparator {
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

			if IsTerminalSymbol(buffer) && leftSideSymbol == nil {
				newSymbol := NewSymbol(buffer)
				leftSideSymbol = &newSymbol
				buffer = ""
				// check if we already this rule and this another alternative
				alts, exists := grammar.impl[newSymbol]
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

func Serialize(grammar Grammar) string {
	var buffer string
	for symbol, alternativesRolls := range grammar.impl {
		for _, roll := range alternativesRolls {
			rule := symbol.ch + " " + ruleSidesSeparator + " "
			for _, symbolInRule := range roll {
				rule += symbolInRule.ch + " "
			}
			rule = strings.Trim(rule, " ")
			buffer += rule + "\n"
		}
	}
	return buffer
}
