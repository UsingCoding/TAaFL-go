package app

import (
	"compiler/pkg/common/grammary"
	"compiler/pkg/common/lexer"
	"fmt"
	"strconv"
	"strings"
)

type AnalyzerErr struct {
	lexem    lexer.Lexem
	expected string
}

func (err AnalyzerErr) Error() string {
	return fmt.Sprintf("Error in %d:%d expected %s got %s", err.lexem.Line, err.lexem.Position, err.expected, err.lexem.Type)
}

func Runner(leftParts, rightParts [][]string, lexerImpl lexer.Lexer) error {
	word, err := lexerImpl.FetchLexem()
	if err != nil {
		return err
	}

	symbol := leftParts[0][0]
	var ptrToRight, ptrToLeft = 0, 0
	var stack []int

	for word.Type != "eof" {
		for i := ptrToLeft; i < len(leftParts); i++ {
			if symbol == leftParts[i][0] {
				if strings.Contains(leftParts[i][1], string(word.Type)) {
					ptrToRight, err = strconv.Atoi(leftParts[i][2])
					if err != nil {
						return err
					}
					break
				} else if leftParts[i][3] == "1" {
					return err
				}
			}
		}

		for i := ptrToRight; i < len(rightParts); i++ {
			if strings.Contains(rightParts[i][1], string(word.Type)) {

				shift, err2 := strconv.Atoi(rightParts[i][3])
				if err2 != nil {
					return err2
				}
				if shift == 1 {
					word, err = lexerImpl.FetchLexem()
					if err != nil {
						return err
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
					if len(stack) > 0 {
						top := len(stack) - 1
						symbol = rightParts[stack[top]][0]
						i = stack[top] - 1
						stack = stack[:top]
						continue
					} else {
						return checkEnd(stack, word, rightParts[i][1])
					}
				}

				if grammary.IsNonTerminalSymbol(rightParts[i][0]) {
					ptrToLeft, err = strconv.Atoi(rightParts[i][2])
					if err != nil {
						return err
					}
					symbol = rightParts[i][0]

					stackPtr, err2 := strconv.Atoi(rightParts[i][4])
					if err2 != nil {
						return err
					}

					if stackPtr > 0 {
						stack = append(stack, stackPtr)
					}
					break
				}

			} else {
				return AnalyzerErr{
					lexem:    word,
					expected: rightParts[i][1],
				}
			}
		}
	}

	return checkEnd(stack, word, string(lexer.LexemTypeEOF))
}

func checkEnd(stack []int, lexem lexer.Lexem, expectedLexem string) error {
	if len(stack) == 0 && lexem.Type == "eof" {
		return nil
	}

	return AnalyzerErr{
		lexem:    lexem,
		expected: expectedLexem,
	}
}
