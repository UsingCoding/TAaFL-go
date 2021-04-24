package app

import (
	"compiler/pkg/common/grammary"
	"strconv"
	"strings"
)

/*
table of left parts of rules:
symbol | NM | ptr to start of right part | error(0/1)
table of right parts of rules:
symbol | NM | ptr to left part of rule or next string | shift(0/1) | stack(0/1)
*/

func CreateTables(data string) ([][]string, [][]string) {
	m := make(map[string]int)
	var leftParts, rightParts [][]string
	var ptrToRightPart int
	for _, line := range strings.Split(data, "\n") {
		elem := ""
		creatingNM, rightPartFlag := false, false
		left, right := []string{}, [][]string{}
		for _, i := range line {
			ch := string(i)
			if ch != " " {
				elem += ch
				continue
			}
			if elem == grammary.RuleSidesSeparator {
				elem = ""
				continue
			}
			if elem == grammary.RuleSequenceSeparator {
				creatingNM = true
				elem = ""
				continue
			}
			if elem == grammary.EmptySymbol {
				temp := []string{}
				temp = append(temp, elem)
				right = append(right, temp)
				elem = ""
				continue
			}
			if !creatingNM {
				if grammary.IsNonTerminalSymbol(elem) {
					if strings.HasPrefix(line, elem) && !rightPartFlag {
						v, ok := m[elem]
						v = 0 + v
						//add nonterminal in map with number of his string
						if !ok {
							m[elem] = len(leftParts) + len(left)
						}
						left = append(left, elem)
						rightPartFlag = true
					} else {
						temp := []string{}
						temp = append(temp, elem)
						right = append(right, temp)
					}
					elem = ""
				} else {
					temp := []string{}
					temp = append(temp, elem)
					temp = append(temp, elem) // add NM to terminal symbol
					right = append(right, temp)
					elem = ""
				}
			}
		}
		left = append(left, elem)                         // add NM to nonterminal
		left = append(left, strconv.Itoa(ptrToRightPart)) // add ptr to right part of rule
		ptrToRightPart += len(right)
		for i := 0; i < len(right); i++ {
			//add NM to nonterminal symbols
			if len(right[i]) < 2 {
				right[i] = append(right[i], elem)
			}
			if right[i][0] == grammary.EmptySymbol {
				right[i] = append(right[i], "nullptr") // add empty ptr field for empty symbol
			}
			// fill ptr to next string int right parts
			if !(grammary.IsNonTerminalSymbol(right[i][0])) {
				if i == (len(right) - 1) {
					if right[i][0] != grammary.EmptySymbol {
						right[i] = append(right[i], "nullptr")
					}
				} else {
					right[i] = append(right[i], strconv.Itoa(len(rightParts)+i+1))
				}
				// fill shift and stack for therminals
				// shift
				if right[i][0] != grammary.EmptySymbol {
					right[i] = append(right[i], "1")
				} else {
					right[i] = append(right[i], "0")
				}
				right[i] = append(right[i], "0") // stack
			} else {
				right[i] = append(right[i], "")  // empty string for ptr to left part of rule
				right[i] = append(right[i], "0") // No shift on nonterminal symbols
				// fill stack field if needed
				if i != (len(right) - 1) {
					right[i] = append(right[i], strconv.Itoa(len(rightParts)+i))
				} else {
					right[i] = append(right[i], "0")
				}
			}
			rightParts = append(rightParts, right[i])
		}
		leftParts = append(leftParts, left)
		creatingNM = false
	}
	// fill error column of left parts
	for i := 0; i < len(leftParts); i++ {
		if i != len(leftParts)-1 {
			if leftParts[i][0] == leftParts[i+1][0] {
				leftParts[i] = append(leftParts[i], strconv.Itoa(0))
			} else {
				leftParts[i] = append(leftParts[i], strconv.Itoa(1))
			}
		} else {
			leftParts[i] = append(leftParts[i], strconv.Itoa(1))
		}
	}
	// add ptr to nonterminal symbols and nullptr to EmptySymbol in right parts
	for i := 0; i < len(rightParts); i++ {
		v, ok := m[rightParts[i][0]]
		if ok {
			rightParts[i][2] = strconv.Itoa(v) // change empty ptr field
		}
	}
	// correct NM in right parts
	for i := 0; i < len(rightParts); i++ {
		if grammary.IsNonTerminalSymbol(rightParts[i][0]) {
			symbol := rightParts[i][0]
			newNM := ""
			ptr, err := strconv.Atoi(rightParts[i][2])
			if err == nil {
				newNM += leftParts[ptr][1]
				ptr += 1
				for ptr < len(leftParts) {
					if leftParts[ptr][0] == symbol {
						newNM += "," + leftParts[ptr][1]
					}
					ptr += 1
				}
			}
			rightParts[i][1] = newNM
		}
	}
	return leftParts, rightParts
}
