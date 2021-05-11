package grammary

import (
	"errors"
	"fmt"
	"strings"
)

func NewSerializer() Serializer {
	return Serializer{}
}

type Serializer struct {
}

func (serializer Serializer) SerializeGrammar(grammar interface{}) (string, error) {
	switch grammar := grammar.(type) {
	case *Grammar:
		return serialize(grammar), nil
	case *GrammarWithHeadSequences:
		return serializeWithSequences(grammar), nil
	default:
		return "", errors.New(fmt.Sprintf("unknown grammar type %T", grammar))
	}

}

func serialize(grammar *Grammar) string {
	var buffer string

	axiomPtr := grammar.Axiom
	buffer += serializeRoll(*axiomPtr, grammar.Impl[*axiomPtr])

	for symbol, alternativesRolls := range grammar.Impl {
		if symbol == *axiomPtr {
			continue
		}
		buffer += serializeRoll(symbol, alternativesRolls)
	}
	return strings.Trim(buffer, "\n")
}

func serializeRoll(leftSideSymbol Symbol, rolls [][]Symbol) string {
	var buffer string
	for _, roll := range rolls {
		rule := leftSideSymbol.ch + " " + RuleSidesSeparator + " "
		for _, symbolInRule := range roll {
			rule += symbolInRule.ch + " "
		}
		rule = strings.Trim(rule, " ")
		buffer += rule + "\n"
	}
	return buffer
}

func serializeWithSequences(grammar *GrammarWithHeadSequences) string {
	var buffer string
	for symbol, alternativesRolls := range grammar.Grammar.Impl {
		for rollNumber, roll := range alternativesRolls {
			rule := symbol.ch + " " + RuleSidesSeparator + " "
			for _, symbolInRule := range roll {
				rule += symbolInRule.ch + " "
			}

			rule += RuleSequenceSeparator + " "
			rule += strings.Join(grammar.Sequences[symbol][rollNumber], " ")
			buffer += rule + "\n"
		}
	}
	return buffer
}
