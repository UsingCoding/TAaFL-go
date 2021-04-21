package main

import (
	"compiler/pkg/common/lexer"
	"fmt"
)

func main() {
	err := runModule()
	if err != nil {
		fmt.Println(err)
	}
}

func runModule() error {
	lexerImpl := newMockLexer()

	for {
		lexem, err := lexerImpl.FetchLexem()
		if err != nil {
			if err == lexer.ErrEOF {
				fmt.Println("program ends")
				return nil
			}
			return err
		}
		if lexem.Type == lexer.LexemTypeError {
			break
		}

		// do business with your lexem
	}

	return nil
}

func newMockLexer() lexer.Lexer {
	return &mockLexer{lexems: []lexer.Lexem{
		{
			Type:     lexer.LexemTypeKeyword,
			Value:    "var",
			Line:     10,
			Position: 20,
		},
		{
			Type:     lexer.LexemTypeError,
			Value:    "kek",
			Line:     10,
			Position: 20,
		},
	}}
}

type mockLexer struct {
	lexems    []lexer.Lexem
	iterCount int
}

func (m *mockLexer) FetchLexem() (lexer.Lexem, error) {
	if m.iterCount == len(m.lexems) {
		return lexer.Lexem{}, lexer.ErrEOF
	}
	return m.lexems[m.iterCount], nil
}
