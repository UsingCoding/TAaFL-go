package app

import "fmt"

type Symbol struct {
	Kind  string
	Value string
}

type Node interface{}

func NewLeafNode(sym Symbol) *LeafNode {
	return &LeafNode{sym: sym}
}

type LeafNode struct {
	Node

	sym Symbol
}

func (n LeafNode) String() string {
	return fmt.Sprintf("Leaf - %s : %s", n.sym.Kind, n.sym.Value)
}

func NewAdditionNode(
	leftOperand StackElementPointer,
	rightOperand StackElementPointer,
) *AdditionNode {
	return &AdditionNode{
		leftOperand:  leftOperand,
		rightOperand: rightOperand,
	}
}

type AdditionNode struct {
	Node

	leftOperand  StackElementPointer
	rightOperand StackElementPointer
}

func (n AdditionNode) String() string {
	return fmt.Sprintf("Add - %d : %d", n.leftOperand, n.rightOperand)
}

func NewSubtractionNode(
	leftOperand StackElementPointer,
	rightOperand StackElementPointer,
) *SubtractionNode {
	return &SubtractionNode{
		leftOperand:  leftOperand,
		rightOperand: rightOperand,
	}
}

type SubtractionNode struct {
	Node

	leftOperand  StackElementPointer
	rightOperand StackElementPointer
}

func (n SubtractionNode) String() string {
	return fmt.Sprintf("sub - %d : %d", n.leftOperand, n.rightOperand)
}
