package app

import (
	"compiler/pkg/common/grammary"
	"compiler/pkg/common/lexer"
	"fmt"
	"strconv"
	"strings"
)

func Runner(leftParts, rightParts [][]string, lexer lexer.Lexer) (bool, error) {
	word, err := lexer.FetchLexem()
	if err != nil {
		return false, err
	}

	symbol := leftParts[0][0]
	var ptrToRight, ptrToLeft = 0, 0
	var stack []int

	for word.Value != "eof" {
		fmt.Println("INF", word.Value)
		for i := ptrToLeft; i < len(leftParts); i++ {
			if symbol == leftParts[i][0] {
				if strings.Contains(leftParts[i][1], word.Value) {
					ptrToRight, err = strconv.Atoi(leftParts[i][2])
					if err != nil {
						return false, err
					}
					break
				} else if leftParts[i][3] == strconv.Itoa(1) {
					return false, nil
				}
			}
		}

		// fmt.Println("Left", symbol, input[wordPos], ptrToLeft, ptrToRight, stack)

		for i := ptrToRight; i < len(rightParts); i++ {
			if strings.Contains(rightParts[i][1], word.Value) {

				shift, err2 := strconv.Atoi(rightParts[i][3])
				if err2 != nil {
					return false, err2
				}
				if shift == 1 {
					word, err = lexer.FetchLexem()
					if err != nil {
						return false, err
					}
				}

				if rightParts[i][0] == grammary.EmptySymbol {
					top := len(stack) - 1
					ptrToRight = stack[top]
					symbol = rightParts[stack[top]][0]
					stack = stack[:top]
					break
				}

				if rightParts[i][2] == "nullptr" && !grammary.IsNonTerminalSymbol(rightParts[i][0]) {
					// fmt.Println("symbol", rightParts[i][0], stack)
					if len(stack) > 0 {
						top := len(stack) - 1
						symbol = rightParts[stack[top]][0]
						i = stack[top] - 1
						stack = stack[:top]
						continue
					} else {
						return checkEnd(stack, word.Value), nil
					}
				}

				if grammary.IsNonTerminalSymbol(rightParts[i][0]) {
					ptrToLeft, err = strconv.Atoi(rightParts[i][2])
					if err != nil {
						return false, err
					}
					symbol = rightParts[i][0]

					stackPtr, err2 := strconv.Atoi(rightParts[i][4])
					if err2 != nil {
						return false, err2
					}

					if stackPtr > 0 {
						stack = append(stack, stackPtr)
					}
					break
				} else {
					fmt.Println("symbol", rightParts[i][0])
				}

			} else {
				return false, nil
			}
		}
	}

	return checkEnd(stack, word.Value), nil
}

func checkEnd(stack []int, lexem string) (answer bool) {
	return len(stack) == 0 && lexem == "eof"
}
