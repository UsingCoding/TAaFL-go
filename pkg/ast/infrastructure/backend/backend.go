package backend

import (
	"encoding/json"
	"io/ioutil"

	"compiler/pkg/ast/infrastructure/executor"
)

const (
	script = ""

	tempFilePath = "/tmp"
	tempFileName = "ast"
)

func NewNodeJSAstBuilderBackend(nodeJSExecutor executor.NodeJS) AstBuilderBackend {
	return &nodeJSAstBuilderBackend{nodeJSExecutor: nodeJSExecutor}
}

type AstBuilderBackend interface {
	Generate(ast json.RawMessage) ([]byte, error)
}

type nodeJSAstBuilderBackend struct {
	nodeJSExecutor executor.NodeJS
}

func (backend *nodeJSAstBuilderBackend) Generate(ast json.RawMessage) ([]byte, error) {
	tempFile, err := ioutil.TempFile(tempFilePath, tempFileName)
	if err != nil {
		return nil, err
	}

	return backend.nodeJSExecutor.Output([]string{tempFile.Name()})
}
