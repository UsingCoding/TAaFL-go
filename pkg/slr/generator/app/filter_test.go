package app

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"compiler/pkg/common/grammary"
	"compiler/pkg/slr/common/inlinedgrammary"
)

func TestEmptySymbolFilter_Filter(t *testing.T) {
	filter := NewEmptySymbolFilter()

	filteredGrammar, err := filter.Filter(inlinedgrammary.New(
		grammary.NewSymbol("<F>"),
		inlinedgrammary.NewRule(grammary.NewSymbol("<F>"), []grammary.Symbol{
			grammary.NewSymbol("<S>"),
		}),
		//inlinedgrammary.NewRule(grammary.NewSymbol("<S>"), []grammary.Symbol{
		//	grammary.NewSymbol("<A>"),
		//	grammary.NewSymbol("d"),
		//	grammary.NewSymbol("<B>"),
		//}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<S>"), []grammary.Symbol{
			grammary.NewSymbol("<A>"),
			grammary.NewSymbol("c"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<B>"), []grammary.Symbol{
			grammary.NewSymbol("<A>"),
			grammary.NewSymbol("b"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<B>"), []grammary.Symbol{
			grammary.NewSymbol("e"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<A>"), []grammary.Symbol{
			grammary.NewSymbol("a"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<A>"), []grammary.Symbol{
			grammary.NewSymbol("e"),
		}),
	))
	assert.NoError(t, err)

	assertGrammarsEquals(t, inlinedgrammary.New(
		grammary.NewSymbol("<F>"),
		inlinedgrammary.NewRule(grammary.NewSymbol("<F>"), []grammary.Symbol{
			grammary.NewSymbol("<S>"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<S>"), []grammary.Symbol{
			grammary.NewSymbol("<A>"),
			grammary.NewSymbol("c"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<B>"), []grammary.Symbol{
			grammary.NewSymbol("<A>"),
			grammary.NewSymbol("b"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<A>"), []grammary.Symbol{
			grammary.NewSymbol("a"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<S>"), []grammary.Symbol{
			grammary.NewSymbol("c"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<B>"), []grammary.Symbol{
			grammary.NewSymbol("b"),
		}),
	),
		filteredGrammar,
	)

}

func assertGrammarsEquals(t *testing.T, expected, actual inlinedgrammary.Grammar) {
	assert.Equal(t, expected.Axiom(), actual.Axiom())
	assert.Equal(t, len(expected.Rules()), len(actual.Rules()))

	for i, expectedRule := range expected.Rules() {
		actualRule := actual.Rules()[i]
		assert.Equal(t, expectedRule, actualRule)
	}
}
