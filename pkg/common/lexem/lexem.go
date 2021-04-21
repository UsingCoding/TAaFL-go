package lexem

type LexemType string

const (
	Keyword LexemType = "keyword"
)

type Lexem struct {
	Type     LexemType
	Value    string
	Line     int
	Position int
}
