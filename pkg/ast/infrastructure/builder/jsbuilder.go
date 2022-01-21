package builder

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"

	"compiler/pkg/ast/app"
	"compiler/pkg/ast/infrastructure/backend"
)

func NewJsBuilder(backend backend.AstBuilderBackend, printAST bool) app.Builder {
	return &jsBuilder{backend: backend}
}

type jsBuilder struct {
	backend  backend.AstBuilderBackend
	printAST bool
}

func (builder *jsBuilder) Build(s app.Stack) (json.RawMessage, error) {
	stack := &stackWrapper{Stack: s}

	strategy := buildStrategy{stack: stack}

	ast, err := strategy.do()
	if err != nil {
		return nil, err
	}

	if builder.printAST {
		_, _ = os.Stdout.WriteString(string(ast))
	}

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

	case app.VariableDeclarationNode:
		nodeStruct := struct {
			NodeType     string `json:"type"`
			Declarations []struct {
				ID struct {
					NodeType string `json:"nodeType"`
					Name     string `json:"name"`
				} `json:"id"`
				Init struct {
					NodeType string `json:"nodeType"`
					Value    string `json:"name"`
				} `json:"init"`
			} `json:"declarations"`
		}{
			NodeType: "VariableDeclaration",
			Declarations: []struct {
				ID struct {
					NodeType string `json:"nodeType"`
					Name     string `json:"name"`
				} `json:"id"`
				Init struct {
					NodeType string `json:"nodeType"`
					Value    string `json:"name"`
				} `json:"init"`
			}{
				{
					ID: struct {
						NodeType string `json:"nodeType"`
						Name     string `json:"name"`
					}{
						NodeType: "Identifier",
						Name:     concreteNode.VarName,
					},
					Init: struct {
						NodeType string `json:"nodeType"`
						Value    string `json:"name"`
					}{
						NodeType: "Literal",
						Value:    concreteNode.Literal,
					},
				},
			},
		}

		return json.Marshal(nodeStruct)

	case app.UnaryExpressionNode:
		nodeStruct := struct {
			NodeType string `json:"nodeType"`
			Operator string `json:"operator"`
			Argument struct {
				NodeType string `json:"nodeType"`
				Value    string `json:"value"`
			} `json:"argument"`
		}{
			NodeType: "UnaryExpression",
			Operator: concreteNode.Operator,
			Argument: struct {
				NodeType string `json:"nodeType"`
				Value    string `json:"value"`
			}{
				NodeType: "Literal",
				Value:    concreteNode.Literal,
			},
		}

		return json.Marshal(nodeStruct)

	case app.VariableOperandNode:
		nodeStruct := struct {
			NodeType string `json:"nodeType"`
			Name     string `json:"name"`
		}{
			NodeType: "Identifier",
			Name:     concreteNode.Name,
		}

		return json.Marshal(nodeStruct)

	case app.AssigmentExpression:
		leftOperandNode, err := s.handleNode(s.stack.get(concreteNode.LeftOperand))
		if err != nil {
			return nil, err
		}
		rightOperandNode, err := s.handleNode(s.stack.get(concreteNode.RightOperand))
		if err != nil {
			return nil, err
		}

		nodeStruct := struct {
			NodeType string          `json:"nodeType"`
			Left     json.RawMessage `json:"left"`
			Right    json.RawMessage `json:"right"`
		}{
			NodeType: "AssignmentExpression",
			Left:     leftOperandNode,
			Right:    rightOperandNode,
		}

		return json.Marshal(nodeStruct)
	case app.ExpressionNode:
		leftOperandNode, err := s.handleNode(s.stack.get(concreteNode.LeftOperand))
		if err != nil {
			return nil, err
		}
		rightOperandNode, err := s.handleNode(s.stack.get(concreteNode.RightOperand))
		if err != nil {
			return nil, err
		}

		nodeStruct := struct {
			NodeType string          `json:"nodeType"`
			Left     json.RawMessage `json:"left"`
			Right    json.RawMessage `json:"right"`
		}{
			NodeType: "ExpressionStatement",
			Left:     leftOperandNode,
			Right:    rightOperandNode,
		}

		return json.Marshal(nodeStruct)

	case app.UpdateExpressionNode:
		argumentNode, err := s.handleNode(s.stack.get(concreteNode.Argument))
		if err != nil {
			return nil, err
		}

		nodeStruct := struct {
			NodeType string          `json:"nodeType"`
			Argument json.RawMessage `json:"argument"`
			Operator string          `json:"operator"`
			Prefix   bool            `json:"prefix"`
		}{
			NodeType: "UpdateExpression",
			Argument: argumentNode,
			Operator: concreteNode.Operator,
			Prefix:   false,
		}

		return json.Marshal(nodeStruct)

	case app.BlockStatement:
		nodeStruct := struct {
			NodeType string            `json:"type"`
			Body     []json.RawMessage `json:"body"`
		}{
			NodeType: "BlockStatement",
		}

		for _, pointer := range concreteNode.Statements {
			operandNode, err := s.handleNode(s.stack.get(pointer))
			if err != nil {
				return nil, err
			}

			nodeStruct.Body = append(nodeStruct.Body, operandNode)
		}

		return json.Marshal(nodeStruct)

	case app.ConditionNode:
		leftOperandNode, err := s.handleNode(s.stack.get(concreteNode.LeftOperand))
		if err != nil {
			return nil, err
		}
		rightOperandNode, err := s.handleNode(s.stack.get(concreteNode.RightOperand))
		if err != nil {
			return nil, err
		}

		nodeStruct := struct {
			NodeType     string          `json:"nodeType"`
			Operator     string          `json:"operator"`
			LeftOperand  json.RawMessage `json:"leftOperand"`
			RightOperand json.RawMessage `json:"rightOperand"`
		}{
			NodeType:     "BinaryExpression",
			Operator:     concreteNode.Operator,
			LeftOperand:  leftOperandNode,
			RightOperand: rightOperandNode,
		}

		return json.Marshal(nodeStruct)

	case app.IFStatement:
		testExpression, err := s.handleNode(s.stack.get(concreteNode.TestExpression))
		if err != nil {
			return nil, err
		}
		consequentExpression, err := s.handleNode(s.stack.get(concreteNode.Consequent))
		if err != nil {
			return nil, err
		}

		var alternateEpxprPtr *json.RawMessage
		if concreteNode.Alternative != nil {
			alternativeExpression, err := s.handleNode(s.stack.get(*concreteNode.Alternative))
			if err != nil {
				return nil, err
			}

			alternateEpxprPtr = &alternativeExpression
		}

		nodeStruct := struct {
			NodeType   string           `json:"nodeType"`
			Test       json.RawMessage  `json:"test"`
			Consequent json.RawMessage  `json:"consequent"`
			Alternate  *json.RawMessage `json:"alternate"`
		}{
			NodeType:   "IfStatement",
			Test:       testExpression,
			Consequent: consequentExpression,
			Alternate:  alternateEpxprPtr,
		}

		return json.Marshal(nodeStruct)
	case app.FORStatement:
		testExpression, err := s.handleNode(s.stack.get(concreteNode.TestExpression))
		if err != nil {
			return nil, err
		}
		varDeclaration, err := s.handleNode(s.stack.get(concreteNode.VariableDeclaration))
		if err != nil {
			return nil, err
		}
		updateExpression, err := s.handleNode(s.stack.get(concreteNode.UpdateExpression))
		if err != nil {
			return nil, err
		}
		body, err := s.handleNode(s.stack.get(concreteNode.Body))
		if err != nil {
			return nil, err
		}

		nodeStruct := struct {
			NodeType         string          `json:"nodeType"`
			Test             json.RawMessage `json:"test"`
			VarDeclaration   json.RawMessage `json:"varDeclaration"`
			UpdateExpression json.RawMessage `json:"updateExpression"`
			Body             json.RawMessage `json:"body"`
		}{
			NodeType:         "ForStatement",
			Test:             testExpression,
			VarDeclaration:   varDeclaration,
			UpdateExpression: updateExpression,
			Body:             body,
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
