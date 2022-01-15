package builder

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"compiler/pkg/ast/app"
	"compiler/pkg/ast/infrastructure/backend"
)

func NewJsBuilder(backend backend.AstBuilderBackend) app.Builder {
	return &jsBuilder{backend: backend}
}

type jsBuilder struct {
	backend backend.AstBuilderBackend
}

func (builder *jsBuilder) Build(s app.Stack) (json.RawMessage, error) {
	stack := &stackWrapper{Stack: s}

	strategy := buildStrategy{stack: stack}

	ast, err := strategy.do()
	if err != nil {
		return nil, err
	}

	fmt.Println(string(ast))

	return nil, errors.New("empty")

	return builder.backend.Generate(ast)
}

type buildStrategy struct {
	stack *stackWrapper
}

func (s *buildStrategy) do() (json.RawMessage, error) {
	topPointer := s.stack.Stack.CurrentTopPointer()
	if topPointer == nil {
		return nil, errors.New("stack has no top")
	}
	top := s.stack.get(*topPointer - 1)

	_, ok := top.(*app.ProgramNode)
	if !ok {
		return nil, errors.Errorf("stack top pointer should be ProgramNode got %T", top)
	}

	program, err := s.handleNode(top)
	if err != nil {
		return nil, err
	}

	return json.Marshal(program)
}

func (s *buildStrategy) handleNode(node app.Node) (json.RawMessage, error) {
	switch concreteNode := node.(type) {
	case *app.ProgramNode:
		nodeStruct := struct {
			NodeType string            `json:"type"`
			Body     []json.RawMessage `json:"body"`
		}{
			NodeType: "Program",
		}

		for _, pointer := range concreteNode.Nodes {
			operandNode, err := s.handleNode(s.stack.get(pointer))
			if err != nil {
				return nil, err
			}

			nodeStruct.Body = append(nodeStruct.Body, operandNode)
		}

		return json.Marshal(nodeStruct)

	case *app.LeafNode:
		nodeStruct := struct {
			NodeType string `json:"type"`
			Value    string `json:"value"`
		}{
			NodeType: "Literal",
			Value:    concreteNode.Sym.Value,
		}

		return json.Marshal(nodeStruct)
	case *app.AdditionNode:
		leftOperandNode, err := s.handleNode(s.stack.get(concreteNode.LeftOperand))
		if err != nil {
			return nil, err
		}
		rightOperandNode, err := s.handleNode(s.stack.get(concreteNode.RightOperand))
		if err != nil {
			return nil, err
		}

		nodeStruct := struct {
			NodeType     string          `json:"type"`
			Operator     string          `json:"operator"`
			LeftOperand  json.RawMessage `json:"left"`
			RightOperand json.RawMessage `json:"right"`
		}{
			NodeType:     "BinaryExpression",
			Operator:     "+",
			LeftOperand:  leftOperandNode,
			RightOperand: rightOperandNode,
		}

		return json.Marshal(nodeStruct)
	case *app.SubtractionNode:
		leftOperandNode, err := s.handleNode(s.stack.get(concreteNode.LeftOperand))
		if err != nil {
			return nil, err
		}
		rightOperandNode, err := s.handleNode(s.stack.get(concreteNode.RightOperand))
		if err != nil {
			return nil, err
		}

		nodeStruct := struct {
			NodeType     string          `json:"type"`
			Operator     string          `json:"operator"`
			LeftOperand  json.RawMessage `json:"left"`
			RightOperand json.RawMessage `json:"right"`
		}{
			NodeType:     "BinaryExpression",
			Operator:     "-",
			LeftOperand:  leftOperandNode,
			RightOperand: rightOperandNode,
		}

		return json.Marshal(nodeStruct)
	}

	return nil, errors.Errorf("unknown node type %T", node)
}
