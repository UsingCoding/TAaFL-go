package app

import (
	"compiler/pkg/common/grammary"
	"fmt"
	"github.com/google/uuid"
)

func FactorizeGrammar(grammar grammary.Grammar) error {
	for leftSideSymbol, altsRolls := range grammar.Impl {
		sameStartedSequenceMap := getRollNumbersStartedWithSameSymbol(altsRolls)
		if len(sameStartedSequenceMap) == 0 {
			continue
		}

		var rollsNumbersToFlush []int

		for symbol, rollsNumbers := range sameStartedSequenceMap {
			newSymbol := fetchNewSymbol()

			grammar.Impl[newSymbol] = populateNewRecord(rollsNumbers, altsRolls)

			for _, rollNumber := range rollsNumbers {
				rollsNumbersToFlush = append(rollsNumbersToFlush, rollNumber)
			}

			grammar.Impl[leftSideSymbol] = append(grammar.Impl[leftSideSymbol], []grammary.Symbol{symbol, newSymbol})
		}

		grammar.Impl[leftSideSymbol] = flushRecords(grammar.Impl[leftSideSymbol], rollsNumbersToFlush)
	}

	return nil
}

func getRollNumbersStartedWithSameSymbol(altsRolls [][]grammary.Symbol) map[grammary.Symbol][]int {
	result := make(map[grammary.Symbol][]int)
	for rollNumber, roll := range altsRolls {
		first := roll[0]

		rollsNumbers, exists := result[first]
		if !exists {
			result[first] = []int{rollNumber}
			continue
		}
		rollsNumbers = append(rollsNumbers, rollNumber)
		result[first] = rollsNumbers
	}

	for symbol, rollNumbers := range result {
		if len(rollNumbers) == 1 {
			delete(result, symbol)
		}
	}

	return result
}

func populateNewRecord(rollsNumbers []int, altsRolls [][]grammary.Symbol) [][]grammary.Symbol {
	var result [][]grammary.Symbol
	for _, rollNumber := range rollsNumbers {
		roll := altsRolls[rollNumber]
		if len(roll) == 1 {
			result = append(result, []grammary.Symbol{grammary.NewSymbol(grammary.EmptySymbol)})
			continue
		}
		result = append(result, roll[1:])
	}
	return result
}

func flushRecords(altsRolls [][]grammary.Symbol, numbersToFlush []int) [][]grammary.Symbol {
	for i, number := range numbersToFlush {
		altsRolls = remove(altsRolls, number-i)
	}
	return altsRolls
}

func fetchNewSymbol() grammary.Symbol {
	return grammary.NewSymbol(fmt.Sprintf("<%s>", uuid.New().String()))
}

func remove(slice [][]grammary.Symbol, s int) [][]grammary.Symbol {
	return append(slice[:s], slice[s+1:]...)
}
