package infrastructure

import (
	"io"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"compiler/pkg/common/lexer"
	"compiler/pkg/lexer/infrastructure/executor"
)

const (
	valueSeparator = "|-|"
	eofSignal      = "eof"
)

func NewLexerAdapter(runtime executor.LexerRuntime) LexerAdapter {
	return &lexerAdapter{
		lexerRuntime: runtime,
	}
}

type LexerAdapter interface {
	lexer.Lexer
	io.Closer
}

type lexerAdapter struct {
	lexerRuntime executor.LexerRuntime
	buffer       []lexer.Lexem
}

func (l *lexerAdapter) FetchLexem() (lexer.Lexem, error) {
	if l.lexerRuntime.IsClosed() {
		return lexer.Lexem{
			Type: lexer.LexemTypeEOF,
		}, nil
	}

	if len(l.buffer) == 0 {
		data, err := l.lexerRuntime.Flush()
		if err != nil {
			return lexer.Lexem{}, err
		}

		buffer, err := parseLexems(data)
		if err != nil {
			return lexer.Lexem{}, err
		}
		l.buffer = buffer
	}

	lexem := l.buffer[0]

	if lexem.Type == lexer.LexemTypeEOF {
		l.lexerRuntime.Close()
	}

	l.buffer = l.buffer[1:]

	return lexem, nil
}

func (l *lexerAdapter) Close() error {
	return l.lexerRuntime.Close()
}

func parseLexems(rowsData string) ([]lexer.Lexem, error) {
	rows := strings.Split(rowsData, "\n")
	var result []lexer.Lexem
	for _, row := range rows {
		if row == "" {
			continue
		}
		lexem, err := parseLexem(row)
		if err != nil {
			return nil, err
		}
		result = append(result, lexem)
	}
	return result, nil
}

func parseLexem(data string) (lexer.Lexem, error) {
	if data == eofSignal {
		return lexer.Lexem{
			Type: lexer.LexemTypeEOF,
		}, nil
	}
	parts := strings.Split(data, valueSeparator)
	if len(parts) != 2 {
		return lexer.Lexem{}, errors.New("unknown lexem format")
	}

	lexemMetaParts := strings.Split(parts[0], " ")
	if len(lexemMetaParts) != 3 {
		return lexer.Lexem{}, errors.New("unknown lexem format in lexem meta info")
	}

	line, err := strconv.Atoi(lexemMetaParts[1])
	if err != nil {
		return lexer.Lexem{}, errors.Wrap(err, "invalid line")
	}

	position, err := strconv.Atoi(lexemMetaParts[2])
	if err != nil {
		return lexer.Lexem{}, errors.Wrap(err, "invalid position")
	}

	lexemType := lexer.LexemType(lexemMetaParts[0])
	//if !isKnownLexemType(lexemType) {
	//	return lexer.Lexem{}, errors.New(fmt.Sprintf("unknown lexem type: '%s'", lexemType))
	//}

	return lexer.Lexem{
		Type:     lexemType,
		Value:    parts[1],
		Line:     line,
		Position: position,
	}, nil
}

func isKnownLexemType(lexemType lexer.LexemType) bool {
	for _, knownLexemType := range lexer.KnownLexemTypes {
		if knownLexemType == lexemType {
			return true
		}
	}
	return false
}
