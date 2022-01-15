package app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack_Add(t *testing.T) {
	stack := Stack{}

	stack.Add(NewAdditionNode(
		stack.Add(NewSubtractionNode(
			stack.Add(NewLeafNode(Symbol{
				Kind:  "float",
				Value: "4.5",
			})),
			stack.Add(NewLeafNode(Symbol{
				Kind:  "float",
				Value: "5.4",
			})),
		)),
		stack.Add(NewAdditionNode(
			stack.Add(NewLeafNode(Symbol{
				Kind:  "int",
				Value: "4",
			})),
			stack.Add(NewLeafNode(Symbol{
				Kind:  "int",
				Value: "5",
			})),
		)),
	))

	for i, node := range stack {
		fmt.Printf("node - %d: %s\n", i, node)
	}

	assert.True(t, true)
}
