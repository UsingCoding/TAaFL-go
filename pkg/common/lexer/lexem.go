package lexer

type LexemType string

const (
	LexemTypeId                 LexemType = "id"
	LexemTypeError              LexemType = "error"
	LexemTypeEOF                LexemType = "eof"
	LexemTypeComparison         LexemType = "comparison"
	LexemTypeSubtraction        LexemType = "subtraction"
	LexemTypeMainToken          LexemType = "mainToken"
	LexemTypeOpenParenthesis    LexemType = "openParenthesis"
	LexemTypeClosingParenthesis LexemType = "closingParenthesis"
	LexemTypeInt                LexemType = "intToken"
	LexemTypeInteger            LexemType = "integer"
	LexemTypeAppropriation      LexemType = "appropriation"
	LexemTypeSeparator          LexemType = "separator"
)

var (
	KnownLexemTypes = []LexemType{
		LexemTypeId,
		LexemTypeError,
		LexemTypeComparison,
		LexemTypeSubtraction,
		LexemTypeMainToken,
		LexemTypeOpenParenthesis,
		LexemTypeClosingParenthesis,
		LexemTypeInt,
		LexemTypeInteger,
		LexemTypeAppropriation,
		LexemTypeSeparator,
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
