package app

import (
	"compiler/pkg/common/grammary"
	"fmt"
	"github.com/google/uuid"
)

func RemoveLeftRecursion(grammar grammary.Grammar) {
	for leftSymbol, altsRolls := range grammar.Impl {
		var repeatedRollsNumbers []int
		for i, roll := range altsRolls {
			if roll[0] == leftSymbol {
				repeatedRollsNumbers = append(repeatedRollsNumbers, i)
			}
		}
		if repeatedRollsNumbers == nil {
			continue
		}

		newNonTerminal := fetchNewNonTerminal()

		// append new non-terminal to other rules in roll
		for i, roll := range altsRolls {
			if inSlice(i, repeatedRollsNumbers) {
				continue
			}

			if len(roll) == 1 && roll[0].String() == grammary.EmptySymbol {
				altsRolls[i] = []grammary.Symbol{newNonTerminal}
				continue
			}

			altsRolls[i] = append(roll, newNonTerminal)
		}

		grammar.Impl[newNonTerminal] = fetchRecordForNewNonTerminal(altsRolls, repeatedRollsNumbers, newNonTerminal)

		grammar.Impl[leftSymbol] = flushRecords(altsRolls, repeatedRollsNumbers)

		break
	}
}

func fetchRecordForNewNonTerminal(altRolls [][]grammary.Symbol, repeatedRollsNumbers []int, newNonTerminal grammary.Symbol) [][]grammary.Symbol {
	// reserve 1 for empty symbol
	result := make([][]grammary.Symbol, len(repeatedRollsNumbers)+1)
	for i, number := range repeatedRollsNumbers {
		result[i] = altRolls[number][1:]

		result[i] = append(result[i], newNonTerminal)
	}
	result[len(repeatedRollsNumbers)] = []grammary.Symbol{grammary.NewSymbol(grammary.EmptySymbol)}
	return result
}

func fetchNewNonTerminal() grammary.Symbol {
	return grammary.NewSymbol(fmt.Sprintf("<%s>", uuid.New().String()))
}

func inSlice(needle int, haystack []int) bool {
	for _, value := range haystack {
		if value == needle {
			return true
		}
	}
	return false
}
