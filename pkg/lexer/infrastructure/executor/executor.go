package executor

import (
	"io"
	"os/exec"
	"strings"
)

type LexerRuntime interface {
	Start() error
	Write(data string) error
	Flush() (string, error)
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
	stdin        *inStream
	stdout       *outStream
}

func (executor *lexerExecutor) Start() error {
	err := executor.lexerProcess.Start()
	if err != nil {
		return err
	}

	executor.stdout.Read()

	return nil
}

func (executor *lexerExecutor) Write(data string) error {
	return executor.stdin.Write(escapeNewLines(data))
}

func (executor *lexerExecutor) Flush() (string, error) {
	err := executor.stdin.Write("! !")
	if err != nil {
		return "", err
	}

	data := executor.stdout.Read()

	return data, nil
}

func (executor *lexerExecutor) Close() error {
	return executor.lexerProcess.Process.Kill()
}

func escapeNewLines(data string) string {
	return strings.Replace(data, "\n", "\\n", -1)
}
