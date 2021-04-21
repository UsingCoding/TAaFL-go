package app

import (
	"compiler/pkg/common/grammary"
	"fmt"
	"github.com/pkg/errors"
)

func BuildHeadSequencesForGrammar(grammar grammary.Grammar) (grammary.GrammarWithHeadSequences, error) {
	sequences := make(map[grammary.Symbol][][]string)

	skippedNonTerminals := make(map[grammary.Symbol]map[int]grammary.Symbol)

	var emptySymbolNonTerminals []grammary.Symbol

	for leftSideSymbol, alternativesRolls := range grammar.Impl {
		for rollNumber, roll := range alternativesRolls {
			symbol := roll[0]
			if symbol.NonTerminal() {
				// we skip populating sequence to fill it with terminals
				symbols, exists := skippedNonTerminals[leftSideSymbol]
				if exists {
					symbols[rollNumber] = symbol
					continue
				}

				skippedNonTerminals[leftSideSymbol] = map[int]grammary.Symbol{rollNumber: symbol}
				continue
			}

			if symbol.String() == grammary.EmptySymbol {
				emptySymbolNonTerminals = append(emptySymbolNonTerminals, leftSideSymbol)
				continue
			}

			altsRolls, ok := sequences[leftSideSymbol]
			if !ok {
				sequences[leftSideSymbol] = [][]string{}
				altsRolls = sequences[leftSideSymbol]
			}

			// verify that alternatives roll exists
			if len(altsRolls) > rollNumber {
				foundRoll := altsRolls[rollNumber]
				foundRoll = append(foundRoll, symbol.String())
				altsRolls[rollNumber] = foundRoll
				sequences[leftSideSymbol] = altsRolls
				continue
			}

			altsRolls = append(altsRolls, []string{symbol.String()})
			sequences[leftSideSymbol] = altsRolls
		}
	}

	if len(sequences) == 0 {
		return grammary.GrammarWithHeadSequences{}, errors.New("no one rules doesnt start with terminal to define head sequence")
	}

	for leftSideSymbol, symbols := range skippedNonTerminals {
		sequence := sequences[leftSideSymbol]
		for rollNumber, symbol := range symbols {
			rolls, exists := sequences[symbol]
			if !exists {
				return grammary.GrammarWithHeadSequences{}, errors.New(fmt.Sprintf("missing nonterminal %s in left side", symbol))
			}

			// copy sequence to append after merging
			afterRollNumberSequencePtr := sequence[rollNumber:]
			afterRollNumberSequence := make([][]string, len(afterRollNumberSequencePtr))
			copy(afterRollNumberSequence, afterRollNumberSequencePtr)

			sequence = append(sequence[:rollNumber], collectHeadingSymbols(rolls))
			sequence = append(sequence, afterRollNumberSequence...)
		}
		sequences[leftSideSymbol] = sequence
	}

	for _, nonTerminal := range emptySymbolNonTerminals {
		foundedSequence, err := findSequenceForNonTerminal(grammar, nonTerminal)
		if err != nil {
			return grammary.GrammarWithHeadSequences{}, err
		}

		sequences[nonTerminal] = [][]string{{foundedSequence}}
	}

	return grammary.GrammarWithHeadSequences{
		Grammar:   grammar,
		Sequences: sequences,
	}, nil
}

func collectHeadingSymbols(sequenceRolls [][]string) []string {
	var res []string
	for _, roll := range sequenceRolls {
		for _, symbol := range roll {
			res = append(res, symbol)
		}
	}
	return res
}

func findSequenceForNonTerminal(grammar grammary.Grammar, nonTerminal grammary.Symbol) (string, error) {
	for leftSideSymbol, altsRolls := range grammar.Impl {
		for _, roll := range altsRolls {
			for i, symbol := range roll {
				if symbol == nonTerminal {
					if len(roll)-1 == i {
						if leftSideSymbol == *grammar.Axiom {
							return grammary.EndOfSequence, nil
						}

						return findSequenceForNonTerminal(grammar, leftSideSymbol)
					} else {
						return roll[i+1].String(), nil
					}
				}
			}
		}
	}

	return "", errors.New(fmt.Sprintf("not found %s in any sequence", nonTerminal))
}
