package lexer

import "errors"

type LexemType string

const (
	LexemTypeKeyword LexemType = "keyword"
	LexemTypeError             = "error"
)

var (
	ErrEOF = errors.New("end of lexems")
)

type Lexem struct {
	Type     LexemType
	Value    string
	Line     int
	Position int
}

type Lexer interface {
	FetchLexem() (Lexem, error)
}
