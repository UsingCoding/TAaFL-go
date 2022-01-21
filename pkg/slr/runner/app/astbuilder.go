package app

import (
	"fmt"

	"github.com/pkg/errors"

	ast "compiler/pkg/ast/app"
)

type astBuilder struct {
	astStack    ast.Stack
	symbolStack *Stack

	pointer *ast.StackElementPointer
}

func (builder astBuilder) buildVariableDeclaration(varName string, varType string, literal string) {
	builder.astStack.Add(ast.NewVariableDeclarationNode(varName, varType, literal))
}

func (builder astBuilder) buildUnaryExpressionOperand(operator string, literal string) {
	builder.astStack.Add(ast.NewUnaryExpressionNode(operator, literal))
}

func (builder astBuilder) buildVariableOperand(varName string) {
	builder.astStack.Add(ast.NewVariableOperandNode(varName))
}

func (builder astBuilder) buildAssigmentExpression() {
	topPointer := builder.astStack.CurrentTopPointer()
	if topPointer == nil {
		panic("nil pointer in ast.Stack while building addition node")
	}
	err := builder.proceedAllowedTypeConversionsByPointers(*topPointer, *topPointer-1)
	if err != nil {
		panic(err)
	}
	builder.astStack.Add(ast.NewAssigmentExpression(*topPointer, *topPointer-1))
}

func (builder astBuilder) buildExpression(operator string) {
	topPointer := builder.astStack.CurrentTopPointer()
	if topPointer == nil {
		panic("nil pointer in ast.Stack while building addition node")
	}
	err := builder.proceedAllowedTypeConversionsByPointers(*topPointer, *topPointer-1)
	if err != nil {
		panic(err)
	}
	builder.astStack.Add(ast.NewExpressionNode(operator, *topPointer, *topPointer-1))
}

func (builder *astBuilder) beginBlockStatement() {
	if builder.pointer != nil {
		panic("astBuilder.pointer already in use")
	}

	builder.pointer = builder.astStack.CurrentTopPointer()
}

func (builder *astBuilder) buildBlockStatement() {
	if builder.pointer == nil {
		panic("astBuilder.pointer not allocated before")
	}
	blockBeginning := *builder.pointer
	builder.pointer = nil

	var pointers []ast.StackElementPointer
	var begin = blockBeginning
	for i := begin; begin <= *builder.astStack.CurrentTopPointer(); i++ {
		pointers = append(pointers, i)
	}

	builder.astStack.Add(ast.NewBlockStatement(pointers))
}

func (builder astBuilder) buildCondition(operator string) {
	topPointer := builder.astStack.CurrentTopPointer()
	if topPointer == nil {
		panic("nil pointer in ast.Stack while building addition node")
	}
	err := builder.proceedAllowedTypeConversionsByPointers(*topPointer, *topPointer-1)
	if err != nil {
		panic(err)
	}
	builder.astStack.Add(ast.NewConditionNode(operator, *topPointer, *topPointer-1))
}

func (builder astBuilder) buildIFStatement() {
	topPointer := builder.astStack.CurrentTopPointer()
	if topPointer == nil {
		panic("nil pointer in ast.Stack while building addition node")
	}

	var alternativePointer *ast.StackElementPointer
	alternative := *topPointer - 2
	_, ok := builder.astStack[alternative].(ast.BlockStatement)
	if !ok {
		// There is no else branch
		alternativePointer = nil
	}

	builder.astStack.Add(ast.NewIFStatement(*topPointer, *topPointer-1, alternativePointer))
}

func (builder astBuilder) updateExpression(operator string) {
	topPointer := builder.astStack.CurrentTopPointer()
	if topPointer == nil {
		panic("nil pointer in ast.Stack while building addition node")
	}

	builder.astStack.Add(ast.NewUpdateExpressionNode(operator, *topPointer))
}

func (builder astBuilder) buildFORStatement() {
	topPointer := builder.astStack.CurrentTopPointer()
	if topPointer == nil {
		panic("nil pointer in ast.Stack while building FOR statement node")
	}

	variableDeclPointer := *topPointer - 3
	testExpressionPointer := *topPointer - 2
	updateExpressionPointer := *topPointer - 1
	bodyPointer := *topPointer

	builder.astStack.Add(ast.NewFORStatement(variableDeclPointer, testExpressionPointer, updateExpressionPointer, bodyPointer))
}

//Math operations
// Don't need any node pointer since all arguments already in astStack
func (builder astBuilder) buildAddition() {
	topPointer := builder.astStack.CurrentTopPointer()
	if topPointer == nil {
		panic("nil pointer in ast.Stack while building addition node")
	}
	err := builder.proceedAllowedTypeConversionsByPointers(*topPointer, *topPointer-1)
	if err != nil {
		panic(err)
	}
	builder.astStack.Add(ast.NewAdditionNode(*topPointer, *topPointer-1))
}

func (builder astBuilder) buildSubtraction() {
	topPointer := builder.astStack.CurrentTopPointer()
	if topPointer == nil {
		panic("nil pointer in ast.Stack while building addition node")
	}
	err := builder.proceedAllowedTypeConversionsByPointers(*topPointer, *topPointer-1)
	if err != nil {
		panic(err)
	}
	builder.astStack.Add(ast.NewSubtractionNode(*topPointer, *topPointer-1))
}

func (builder astBuilder) buildMultiplication() {
	topPointer := builder.astStack.CurrentTopPointer()
	if topPointer == nil {
		panic("nil pointer in ast.Stack while building addition node")
	}
	err := builder.proceedAllowedTypeConversionsByPointers(*topPointer, *topPointer-1)
	if err != nil {
		panic(err)
	}
	builder.astStack.Add(ast.NewMultiplicationNode(*topPointer, *topPointer-1))
}

func (builder astBuilder) buildDivision() {
	topPointer := builder.astStack.CurrentTopPointer()
	if topPointer == nil {
		panic("nil pointer in ast.Stack while building addition node")
	}
	err := builder.proceedAllowedTypeConversionsByPointers(*topPointer, *topPointer-1)
	if err != nil {
		panic(err)
	}
	builder.astStack.Add(ast.NewDivisionNode(*topPointer, *topPointer-1))
}

func (builder astBuilder) proceedAllowedTypeConversions(firstType, secondType string) error {
	allower := func(firstType, secondType string) error {
		if firstType == "boolean" && secondType == "string" {
			return errors.Errorf("cannot convert properly type %s to type %s in condition statement", firstType, secondType)
		}

		if firstType == "boolean" && secondType == "int" {
			return errors.Errorf("cannot convert properly type %s to type %s in condition statement", firstType, secondType)
		}

		if firstType == "boolean" && secondType == "double" {
			return errors.Errorf("cannot convert properly type %s to type %s in condition statement", firstType, secondType)
		}

		return nil
	}

	err := allower(firstType, secondType)
	if err != nil {
		return err
	}

	return allower(secondType, firstType)
}

func (builder astBuilder) proceedAllowedTypeConversionsByPointers(firstPointer, secondPointer ast.StackElementPointer) error {
	return builder.proceedAllowedTypeConversions(
		builder.fetchTypeFromNode(builder.astStack[firstPointer]),
		builder.fetchTypeFromNode(builder.astStack[secondPointer]),
	)
}

func (builder astBuilder) fetchTypeFromNode(node ast.Node) string {
	switch concreteNode := node.(type) {
	case ast.LeafNode:
		return concreteNode.Sym.Kind
	case ast.VariableOperandNode:
		varName := concreteNode.Name
		variable := builder.symbolStack.GetLast().Find(varName)
		if variable == nil {
			panic(fmt.Sprintf("cannot find name %s", varName))
		}

		return variable.kind
	}

	panic("unknown node type")
}
