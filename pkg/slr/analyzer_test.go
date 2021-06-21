package slr

import (
	"compiler/pkg/common/grammary"
	commonlexer "compiler/pkg/common/lexer"
	"compiler/pkg/slr/common"
	"compiler/pkg/slr/common/inlinedgrammary"
	slrgenerator "compiler/pkg/slr/generator/app"
	slrrunnner "compiler/pkg/slr/runner/app"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAnalyzer_Analyze1(t *testing.T) {
	// we put no lexems cause we have error before runner will proceed sequnce
	analyzer := buildAnalyzer(commonlexer.Lexem{
		Type:  commonlexer.LexemTypeOpenParenthesis,
		Value: "(",
	},
		commonlexer.Lexem{
			Type:  commonlexer.LexemTypeOpenParenthesis,
			Value: "(",
		},
		commonlexer.Lexem{
			Type:  commonlexer.LexemTypeClosingParenthesis,
			Value: ")",
		},
		commonlexer.Lexem{
			Type:  commonlexer.LexemTypeClosingParenthesis,
			Value: ")",
		},
		commonlexer.Lexem{
			Type:  commonlexer.LexemTypeEOF,
			Value: "_|_",
		})

	err := analyzer.Analyze(inlinedgrammary.New(
		grammary.NewSymbol("<F>"),
		inlinedgrammary.NewRule(grammary.NewSymbol("<F>"), []grammary.Symbol{
			grammary.NewSymbol("<S>"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<S>"), []grammary.Symbol{
			grammary.NewSymbol("("),
			grammary.NewSymbol("<S>"),
			grammary.NewSymbol(")"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<S>"), []grammary.Symbol{
			grammary.NewSymbol("("),
			grammary.NewSymbol(")"),
		}),
	))

	assert.NoError(t, err)
}

func TestAnalyzer_Analyze2(t *testing.T) {
	// we put no lexems cause we have error before runner will proceed sequnce
	analyzer := buildAnalyzer()

	err := analyzer.Analyze(inlinedgrammary.New(
		grammary.NewSymbol("<F>"),
		inlinedgrammary.NewRule(grammary.NewSymbol("<F>"), []grammary.Symbol{
			grammary.NewSymbol("<E>"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<E>"), []grammary.Symbol{
			grammary.NewSymbol("<E>"),
			grammary.NewSymbol("+"),
			grammary.NewSymbol("<E>"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<E>"), []grammary.Symbol{
			grammary.NewSymbol("i"),
		}),
	))

	assert.Error(t, err, slrgenerator.ErrNotSlrGrammar{})
}

func buildAnalyzer(lexems ...commonlexer.Lexem) Analyzer {
	generator := slrgenerator.NewGenerator()
	validator := slrgenerator.NewValidator()
	runner := slrrunnner.NewRunner(&mockLexer{lexems: lexems})

	return NewAnalyzer(
		generator,
		validator,
		&mockExporter{},
		runner,
	)
}

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

type mockExporter struct {
}

func (m *mockExporter) Export(common.Table) error {
	return nil
}
