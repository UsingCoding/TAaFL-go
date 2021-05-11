package infrastructure

import (
	"compiler/pkg/common/lexer"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLexerAdapter_FetchLexem(t *testing.T) {
	adapter := NewLexerAdapter(&mockLexerRuntime{rows: []string{
		fmt.Sprintf("id 1 4%slexem content\nid 2 4%slexem content\n", valueSeparator, valueSeparator),
		fmt.Sprintf("id 5 2%slexem", valueSeparator),
	}})

	{
		lexem, err := adapter.FetchLexem()
		assert.NoError(t, err)

		assert.Equal(t, lexer.LexemTypeId, lexem.Type)
		assert.Equal(t, 1, lexem.Line)
		assert.Equal(t, 4, lexem.Position)
		assert.Equal(t, "lexem content", lexem.Value)
	}

	{
		lexem, err := adapter.FetchLexem()
		assert.NoError(t, err)

		assert.Equal(t, lexer.LexemTypeId, lexem.Type)
		assert.Equal(t, 2, lexem.Line)
		assert.Equal(t, 4, lexem.Position)
		assert.Equal(t, "lexem content", lexem.Value)
	}

	{
		lexem, err := adapter.FetchLexem()
		assert.NoError(t, err)

		assert.Equal(t, lexer.LexemTypeId, lexem.Type)
		assert.Equal(t, 5, lexem.Line)
		assert.Equal(t, 2, lexem.Position)
		assert.Equal(t, "lexem", lexem.Value)
	}

	{
		lexem, err := adapter.FetchLexem()
		assert.NoError(t, err)

		assert.Equal(t, lexer.LexemTypeEOF, lexem.Type)
	}

	{
		lexem, err := adapter.FetchLexem()
		assert.NoError(t, err)

		assert.Equal(t, lexer.LexemTypeEOF, lexem.Type)
	}
}

type mockLexerRuntime struct {
	rows []string
}

func (m *mockLexerRuntime) Start() error {
	return nil
}

func (m *mockLexerRuntime) Write(data string) error {
	return nil
}

func (m *mockLexerRuntime) Flush() (string, error) {
	if len(m.rows) == 0 {
		panic("access to closed runtime")
	}
	row := m.rows[0]
	m.rows = m.rows[1:]
	return row, nil
}

func (m *mockLexerRuntime) Close() error {
	return nil
}

func (m *mockLexerRuntime) IsClosed() bool {
	return len(m.rows) == 0
}
