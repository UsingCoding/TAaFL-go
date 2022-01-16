package backend

import (
	"bytes"
	"encoding/json"

	"compiler/pkg/ast/infrastructure/executor"
)

func NewJSASTBuilderBackend(backendExecutor executor.ASTBackendExecutor) AstBuilderBackend {
	return &nodeJSAstBuilderBackend{backendExecutor: backendExecutor}
}

type AstBuilderBackend interface {
	Generate(ast json.RawMessage) ([]byte, error)
}

type nodeJSAstBuilderBackend struct {
	backendExecutor executor.ASTBackendExecutor
}

func (backend *nodeJSAstBuilderBackend) Generate(ast json.RawMessage) ([]byte, error) {
	return backend.backendExecutor.Output([]string{}, executor.WithStdin(bytes.NewBuffer(ast)))
}
