package executor

import (
	"github.com/pkg/errors"
	"io"
	"os/exec"
	"strings"
)

type LexerRuntime interface {
	Start() error
	Write(data string) error
	Flush() (string, error)
	IsClosed() bool
	io.Closer
}

func NewLexerExecutor(lexerExecutablePath string) LexerRuntime {
	stdin := newInStream()
	stdout := newOutStream()

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
		stdin:        stdin,
		stdout:       stdout,
	}
}

type lexerExecutor struct {
	lexerProcess exec.Cmd
	isClosed     bool
	stdin        *inStream
	stdout       *outStream
}

func (executor *lexerExecutor) Start() error {
	err := executor.lexerProcess.Start()
	if err != nil {
		return newLexerRuntimeError(err)
	}

	executor.stdout.Read()

	return nil
}

func (executor *lexerExecutor) Write(data string) error {
	return newLexerRuntimeError(executor.stdin.Write(escapeNewLines(data)))
}

func (executor *lexerExecutor) Flush() (string, error) {
	if executor.isClosed {
		return "", errors.New("flush in closed runtime")
	}
	err := executor.stdin.Write("! !")
	if err != nil {
		return "", newLexerRuntimeError(err)
	}

	data := executor.stdout.Read()

	return data, nil
}

func (executor *lexerExecutor) Close() error {
	executor.isClosed = true
	return newLexerRuntimeError(executor.lexerProcess.Process.Kill())
}

func (executor *lexerExecutor) IsClosed() bool {
	return executor.isClosed
}

func escapeNewLines(data string) string {
	return strings.Replace(data, "\n", "\\n", -1)
}
