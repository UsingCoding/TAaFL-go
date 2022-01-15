package grammary

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

func Parse(rawData string) (Grammar, error) {
	grammar := NewGrammar()

	for lineNumber, line := range strings.Split(rawData, "\n") {
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
					fmt.Println(buffer)
					return Grammar{}, errors.New(fmt.Sprintf("no leftSideSymbol found on line %d", lineNumber+1))
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
