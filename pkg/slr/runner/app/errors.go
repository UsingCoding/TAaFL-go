package app

import (
	"fmt"

	"compiler/pkg/common/grammary"
	"compiler/pkg/common/lexer"
	"compiler/pkg/slr/common"
)

type ErrUnexpectedSymbol struct {
	expectedSymbols  []grammary.Symbol
	unexpectedSymbol grammary.Symbol
	lexem            lexer.Lexem
}

func (e *ErrUnexpectedSymbol) Error() string {
	return fmt.Sprintf(
		"unexpected %s at line %d position %d: expected one of %s",
		e.unexpectedSymbol,
		e.lexem.Line,
		e.lexem.Position,
		e.expectedSymbols,
	)
}

type ErrUnexpectedEnd struct {
	expectedSymbols []grammary.Symbol
}

func (e *ErrUnexpectedEnd) Error() string {
	return fmt.Sprintf(
		"unexpected end: expected one of %s",
		e.expectedSymbols,
	)
}

func (strategy *proceedStrategy) prepareUnexpectedSymbolErr(
	currentState common.TableRef,
	unexpectedSymbol grammary.Symbol,
	lexem lexer.Lexem,
) error {
	var expectedSymbols []grammary.Symbol
	for symbol := range strategy.table.TableMap[currentState] {
		const appendNonTerminals = true

		// Append only terminals
		// Intentionally append nonTerminal for debug purposes
		if !symbol.NonTerminal() || appendNonTerminals {
			expectedSymbols = append(expectedSymbols, symbol)
		}
	}

	return &ErrUnexpectedSymbol{
		expectedSymbols:  expectedSymbols,
		unexpectedSymbol: unexpectedSymbol,
		lexem:            lexem,
	}
}

func (strategy *proceedStrategy) prepareUnexpectedEndErr(currentState common.TableRef) error {
	var expectedSymbols []grammary.Symbol
	for symbol := range strategy.table.TableMap[currentState] {
		expectedSymbols = append(expectedSymbols, symbol)
	}

	return &ErrUnexpectedEnd{
		expectedSymbols: expectedSymbols,
	}
}
