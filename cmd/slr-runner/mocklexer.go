package main

import (
	commonlexer "compiler/pkg/common/lexer"
	"github.com/pkg/errors"
)

type mockLexer struct {
	lexems []commonlexer.Lexem
}

func (lexer *mockLexer) FetchLexem() (commonlexer.Lexem, error) {
	if len(lexer.lexems) == 0 {
		return commonlexer.Lexem{}, errors.New("end of input")
	}
	lexem := lexer.lexems[0]
	lexer.lexems = lexer.lexems[1:]
	return lexem, nil
}
