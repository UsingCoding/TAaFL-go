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
	for symbol, alternativesRolls := range grammar.Impl {
		for _, roll := range alternativesRolls {
			rule := symbol.ch + " " + ruleSidesSeparator + " "
			for _, symbolInRule := range roll {
				rule += symbolInRule.ch + " "
			}
			rule = strings.Trim(rule, " ")
			buffer += rule + "\n"
		}
	}
	return strings.Trim(buffer, "\n")
}

func serializeWithSequences(grammar *GrammarWithHeadSequences) string {
	var buffer string
	for symbol, alternativesRolls := range grammar.Grammar.Impl {
		for rollNumber, roll := range alternativesRolls {
			rule := symbol.ch + " " + ruleSidesSeparator + " "
			for _, symbolInRule := range roll {
				rule += symbolInRule.ch + " "
			}

			rule += ruleSequenceSeparator + " "
			rule += strings.Join(grammar.Sequences[symbol][rollNumber], " ")
			buffer += rule + "\n"
		}
	}
	return buffer
}
