package infrastructure

import (
	"compiler/pkg/common/lexer"
	"io"
)

func NewLexerAdapter() lexer.Lexer {
	return &lexerAdapter{}
}

type LexerAdapter interface {
	lexer.Lexer
	io.Closer
}

type lexerAdapter struct {
	lexerRuntime LexerRuntime
}

func (l *lexerAdapter) FetchLexem() (lexer.Lexem, error) {
	data, err := l.lexerRuntime.Flush()
	if err != nil {
		return lexer.Lexem{}, err
	}

}

func (l *lexerAdapter) Close() error {
	return l.lexerRuntime.Close()
}

func parseLexem(data string) (lexer.Lexem, error) {

}
