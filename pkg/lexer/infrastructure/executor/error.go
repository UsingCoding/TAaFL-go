package executor

func newLexerRuntimeError(err error) error {
	if err == nil {
		return nil
	}
	return LexerRuntimeError{err: err}
}

type LexerRuntimeError struct {
	err error
}

func (l LexerRuntimeError) Error() string {
	const prefix = "LexerRuntimeError: "
	return prefix + l.err.Error()
}
