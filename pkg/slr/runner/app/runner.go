package app

import (
	ast "compiler/pkg/ast/app"
	"compiler/pkg/common/grammary"
	"compiler/pkg/common/lexer"
	"compiler/pkg/slr/common"
)

type Runner interface {
	Proceed(table common.Table, axiom grammary.Symbol) (ast.Stack, error)
}

func NewRunner(lexer lexer.Lexer) Runner {
	return &runner{lexer: lexer}
}

type runner struct {
	lexer lexer.Lexer
}

func (runner *runner) Proceed(table common.Table, axiom grammary.Symbol) (ast.Stack, error) {
	strategy := proceedStrategy{
		axiom: axiom,
		input: &inputStream{lexer: runner.lexer},
		table: table,
	}

	return strategy.do()
}

type proceedStrategy struct {
	axiom      grammary.Symbol
	input      inputLexemStream
	table      common.Table
	stateStack []common.TableRef
}

func (strategy *proceedStrategy) do() (ast.Stack, error) {
	// add axiom table ref as bottom of stack
	strategy.stateStack = append(strategy.stateStack, strategy.table.AxiomRef)

	for {
		lexem, err := strategy.input.fetch()
		if err != nil {
			return nil, err
		}

		if lexem.Type == lexer.LexemTypeEOF && len(strategy.stateStack) == 1 {
			return nil, strategy.prepareUnexpectedEndErr(strategy.stateStack[0])
		}

		lexemValue := lexem.Value
		state := strategy.stateStack[len(strategy.stateStack)-1]

		expectedSymbol := grammary.NewSymbol(lexemValue)

		if state == strategy.table.AxiomRef && expectedSymbol == strategy.axiom {
			// input is OK
			return ast.Stack{}, nil
		}

		tableRef, exists := strategy.table.TableMap[state][expectedSymbol]
		if !exists {
			return nil, strategy.prepareUnexpectedSymbolErr(state, expectedSymbol, lexem)
		}

		tableEntry := strategy.table.TableRefs[tableRef]

		if tableEntry.CollapseEntry == nil {
			strategy.putStateToStack(tableRef)
			continue
		}

		strategy.proceedCollapse(tableEntry, expectedSymbol)
	}
}

func (strategy *proceedStrategy) putStateToStack(tableRef common.TableRef) {
	strategy.stateStack = append(strategy.stateStack, tableRef)
}

func (strategy *proceedStrategy) proceedCollapse(
	tableEntry common.TableEntry,
	expectedSymbol grammary.Symbol,
) {
	collapseEntry := tableEntry.CollapseEntry

	if collapseEntry.CountOfSymbolsInRule == uint(len(strategy.stateStack)) {
		panic("you cannot clear all stack : err in generator")
	}

	strategy.stateStack = strategy.stateStack[:uint(len(strategy.stateStack))-collapseEntry.CountOfSymbolsInRule]

	strategy.input.prepend(expectedSymbol)
	strategy.input.prepend(collapseEntry.Symbol)
}
