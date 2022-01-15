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
