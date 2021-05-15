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
		return serializeWithSequence(grammar), nil
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

func serializeWithSequence(grammar *GrammarWithHeadSequences) string {
	var buffer string

	axiomPtr := grammar.Grammar.Axiom
	buffer += serializeRollWithHeadSequence(*axiomPtr, grammar.Grammar.Impl[*axiomPtr], grammar.Sequences[*axiomPtr])

	for symbol, alternativesRolls := range grammar.Grammar.Impl {
		if symbol == *axiomPtr {
			continue
		}
		buffer += serializeRollWithHeadSequence(symbol, alternativesRolls, grammar.Sequences[symbol])
	}
	return strings.Trim(buffer, "\n")
}

func serializeRollWithHeadSequence(leftSideSymbol Symbol, rolls [][]Symbol, sequence [][]string) string {
	var buffer string
	for rollNumber, roll := range rolls {
		rule := leftSideSymbol.ch + " " + RuleSidesSeparator + " "
		for _, symbolInRule := range roll {
			rule += symbolInRule.ch + " "
		}
		rule += RuleSequenceSeparator + " "
		rule += strings.Join(sequence[rollNumber], " ")
		buffer += rule + "\n"
	}
	return buffer
}
