package app

import (
	"fmt"

	"github.com/pkg/errors"

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
	symbolStack := &Stack{}
	strategy := proceedStrategy{
		axiom:       axiom,
		input:       &inputStream{lexer: runner.lexer},
		table:       table,
		symbolTable: symbolStack,
		astBuilder: &astBuilder{
			symbolStack: symbolStack,
		},
	}

	return strategy.do()
}

type proceedStrategy struct {
	axiom      grammary.Symbol
	input      inputLexemStream
	table      common.Table
	stateStack []common.TableRef

	astBuilder  *astBuilder
	symbolTable *Stack
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

		err = strategy.proceedCollapse(tableEntry, expectedSymbol)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("at %d line and %d row", lexem.Line, lexem.Position))
		}
	}
}

func (strategy *proceedStrategy) putStateToStack(tableRef common.TableRef) {
	strategy.stateStack = append(strategy.stateStack, tableRef)
}

func (strategy *proceedStrategy) proceedCollapse(
	tableEntry common.TableEntry,
	expectedSymbol grammary.Symbol,
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	collapseEntry := tableEntry.CollapseEntry

	if collapseEntry.CountOfSymbolsInRule == uint(len(strategy.stateStack)) {
		panic("you cannot clear all stack : err in generator")
	}

	strategy.stateStack = strategy.stateStack[:uint(len(strategy.stateStack))-collapseEntry.CountOfSymbolsInRule]

	strategy.input.prepend(expectedSymbol)
	strategy.input.prepend(collapseEntry.Symbol)

	switch tableEntry.Rule {
	case common.VariableDeclaration:
		strategy.astBuilder.buildVariableDeclaration(
			strategy.input.getFromTop(4).Value, // Name
			strategy.input.getFromTop(3).Value, // Type
			strategy.input.getFromTop(1).Value, // Value
		)
		sym := Symbol{
			kind: strategy.input.getFromTop(3).Value,
			name: strategy.input.getFromTop(4).Value,
		}
		strategy.symbolTable.GetLast().AddSymbol(sym)
	case common.UnaryExpression:
		strategy.astBuilder.buildUnaryExpressionOperand(
			strategy.input.getFromTop(1).Value, // Operand
			strategy.input.getFromTop(2).Value, // Literal
		)
	case common.VariableOperand:
		varName := strategy.input.getFromTop(1).Value
		if variable := strategy.symbolTable.GetLast().Find(varName); variable == nil {
			return errors.Errorf("cannot find name %s", varName)
		}
		strategy.astBuilder.buildVariableOperand(varName) // Only variable name
	case common.AssigmentExpression:
		strategy.astBuilder.buildAssigmentExpression()
	case common.Expression:
		strategy.astBuilder.buildExpression(strategy.input.getFromTop(1).Value)
	case common.BeginBlockStatement:
		strategy.symbolTable.CreateStack()
		strategy.astBuilder.beginBlockStatement()
	case common.BlockStatement:
		strategy.astBuilder.buildBlockStatement()
		strategy.symbolTable.DeleteLast()
	case common.Condition:
		strategy.astBuilder.buildCondition(strategy.input.getFromTop(1).Value)
	case common.IFStatement:
		strategy.astBuilder.buildIFStatement()
		strategy.symbolTable.DeleteLast()
	case common.UpdateExpression:
		strategy.astBuilder.updateExpression(strategy.input.getFromTop(2).Value)
	case common.FORStatement:
		strategy.astBuilder.buildFORStatement()

	//	Math operations
	case common.Addition:
		strategy.astBuilder.buildAddition()
	case common.Subtraction:
		strategy.astBuilder.buildSubtraction()
	case common.Multiplication:
		strategy.astBuilder.buildMultiplication()
	case common.Division:
		strategy.astBuilder.buildDivision()
	}

	return nil
}
