package app

import "fmt"

type Symbol struct {
	Kind  string
	Value string
}

type Node interface{}

func NewLeafNode(sym Symbol) *LeafNode {
	return &LeafNode{Sym: sym}
}

type LeafNode struct {
	Node

	Sym Symbol
}

func (n LeafNode) String() string {
	return fmt.Sprintf("Leaf - %s : %s", n.Sym.Kind, n.Sym.Value)
}

func NewProgramNode(nodes []StackElementPointer) *ProgramNode {
	return &ProgramNode{Nodes: nodes}
}

// ProgramNode is Node grouping nodes to program
type ProgramNode struct {
	Node

	Nodes []StackElementPointer
}

func (n ProgramNode) String() string {
	return fmt.Sprintf("Program - %d nodes", len(n.Nodes))
}

func NewVariableDeclarationNode(varName string, varType string, literal string) *VariableDeclarationNode {
	return &VariableDeclarationNode{
		VarName: varName,
		VarType: varType,
		Literal: literal,
	}
}

type VariableDeclarationNode struct {
	Node

	VarName string
	VarType string
	Literal string
}

func NewUnaryExpressionNode(operator string, literal string) *UnaryExpressionNode {
	return &UnaryExpressionNode{Operator: operator, Literal: literal}
}

type UnaryExpressionNode struct {
	Node

	Operator string
	Literal  string
}

func NewVariableOperandNode(name string) *VariableOperandNode {
	return &VariableOperandNode{Name: name}
}

type VariableOperandNode struct {
	Node

	Name string
}

func NewAssigmentExpression(leftOperand StackElementPointer, rightOperand StackElementPointer) *AssigmentExpression {
	return &AssigmentExpression{LeftOperand: leftOperand, RightOperand: rightOperand}
}

type AssigmentExpression struct {
	Node

	LeftOperand  StackElementPointer
	RightOperand StackElementPointer
}

func NewExpressionNode(operator string, leftOperand StackElementPointer, rightOperand StackElementPointer) *ExpressionNode {
	return &ExpressionNode{Operator: operator, LeftOperand: leftOperand, RightOperand: rightOperand}
}

type ExpressionNode struct {
	Node

	Operator     string
	LeftOperand  StackElementPointer
	RightOperand StackElementPointer
}

func NewBlockStatement(statements []StackElementPointer) *BlockStatement {
	return &BlockStatement{Statements: statements}
}

type BlockStatement struct {
	Node

	Statements []StackElementPointer
}

func NewConditionNode(operator string, leftOperand StackElementPointer, rightOperand StackElementPointer) *ConditionNode {
	return &ConditionNode{Operator: operator, LeftOperand: leftOperand, RightOperand: rightOperand}
}

type ConditionNode struct {
	Node

	Operator     string
	LeftOperand  StackElementPointer
	RightOperand StackElementPointer
}

func NewIFStatement(testExpression StackElementPointer, consequent StackElementPointer, alternative *StackElementPointer) *IFStatement {
	return &IFStatement{TestExpression: testExpression, Consequent: consequent, Alternative: alternative}
}

type IFStatement struct {
	Node

	TestExpression StackElementPointer
	Consequent     StackElementPointer
	Alternative    *StackElementPointer
}

func NewUpdateExpressionNode(operator string, argument StackElementPointer) *UpdateExpressionNode {
	return &UpdateExpressionNode{Operator: operator, Argument: argument}
}

type UpdateExpressionNode struct {
	Node

	Operator string
	Argument StackElementPointer
}

func NewFORStatement(
	variableDeclaration StackElementPointer,
	testExpression StackElementPointer,
	updateExpression StackElementPointer,
	body StackElementPointer,
) *FORStatement {
	return &FORStatement{VariableDeclaration: variableDeclaration, TestExpression: testExpression, UpdateExpression: updateExpression, Body: body}
}

type FORStatement struct {
	Node

	VariableDeclaration StackElementPointer
	TestExpression      StackElementPointer
	UpdateExpression    StackElementPointer // i++ for i.e
	Body                StackElementPointer
}

// Math operations
func NewAdditionNode(
	leftOperand StackElementPointer,
	rightOperand StackElementPointer,
) *AdditionNode {
	return &AdditionNode{
		LeftOperand:  leftOperand,
		RightOperand: rightOperand,
	}
}

type AdditionNode struct {
	Node

	LeftOperand  StackElementPointer
	RightOperand StackElementPointer
}

func (n AdditionNode) String() string {
	return fmt.Sprintf("Add - %d : %d", n.LeftOperand, n.RightOperand)
}

func NewSubtractionNode(
	leftOperand StackElementPointer,
	rightOperand StackElementPointer,
) *SubtractionNode {
	return &SubtractionNode{
		LeftOperand:  leftOperand,
		RightOperand: rightOperand,
	}
}

type SubtractionNode struct {
	Node

	LeftOperand  StackElementPointer
	RightOperand StackElementPointer
}

func (n SubtractionNode) String() string {
	return fmt.Sprintf("sub - %d : %d", n.LeftOperand, n.RightOperand)
}

func NewMultiplicationNode(leftOperand StackElementPointer, rightOperand StackElementPointer) *MultiplicationNode {
	return &MultiplicationNode{LeftOperand: leftOperand, RightOperand: rightOperand}
}

type MultiplicationNode struct {
	Node

	LeftOperand  StackElementPointer
	RightOperand StackElementPointer
}

func NewDivisionNode(leftOperand StackElementPointer, rightOperand StackElementPointer) *DivisionNode {
	return &DivisionNode{LeftOperand: leftOperand, RightOperand: rightOperand}
}

type DivisionNode struct {
	Node

	LeftOperand  StackElementPointer
	RightOperand StackElementPointer
}
