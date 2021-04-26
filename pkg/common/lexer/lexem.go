package lexer

type LexemType string

const (
	LexemTypeKeyword LexemType = "keyword"
	LexemTypeInteger LexemType = "integer"
	LexemTypeFloat   LexemType = "float"
	LexemTypeError   LexemType = "error"
	LexemTypeEOF     LexemType = "eof"
)

var (
	KnownLexemTypes = []LexemType{
		LexemTypeKeyword,
		LexemTypeInteger,
		LexemTypeFloat,
		LexemTypeError,
	}
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
