package app

import (
	"compiler/pkg/common/grammary"
	"compiler/pkg/common/lexer"
)

type inputLexemStream interface {
	fetch() (lexer.Lexem, error)
	prepend(symbol grammary.Symbol)
	getFromTop(pointer int) lexer.Lexem
}

type inputStream struct {
	lexer  lexer.Lexer
	buffer []lexer.Lexem
}

func (input *inputStream) prepend(symbol grammary.Symbol) {
	input.buffer = append([]lexer.Lexem{{Value: symbol.String()}}, input.buffer...)
}

func (input *inputStream) fetch() (lexer.Lexem, error) {
	if len(input.buffer) != 0 {
		lexem := input.buffer[0]
		input.buffer = input.buffer[1:]
		return lexem, nil
	}

	return input.lexer.FetchLexem()
}

func (input *inputStream) getFromTop(pointer int) lexer.Lexem {
	return input.buffer[len(input.buffer)-pointer]
}
