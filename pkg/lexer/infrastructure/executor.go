package infrastructure

import (
	"bytes"
	"io"
	"os/exec"
)

const (
	lexerExecutablePath = "lexer"
)

type LexerRuntime interface {
	Start() error
	Flush() (string, error)
	io.Closer
}

func NewLexerExecutor() LexerRuntime {
	stdin := bytes.Buffer{}
	stdout := bytes.Buffer{}
	cmd := exec.Cmd{
		Path:         lexerExecutablePath,
		Stdin:        stdin,
		Stdout:       stdout,
		SysProcAttr:  nil,
		Process:      nil,
		ProcessState: nil,
	}

	return &lexerExecutor{
		lexerProcess: cmd,
	}
}

type lexerExecutor struct {
	lexerProcess exec.Cmd
	readChannel  chan string
	writeChannel chan string
}

func (executor *lexerExecutor) Start() error {
	return executor.lexerProcess.Start()
}

func (executor *lexerExecutor) Flush() (string, error) {

}

func (executor *lexerExecutor) Close() error {
	return executor.lexerProcess.Process.Kill()
}
