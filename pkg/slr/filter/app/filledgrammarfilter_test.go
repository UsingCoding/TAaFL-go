package app

import (
	"compiler/pkg/common/grammary"
	"compiler/pkg/slr/common/inlinedgrammary"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilledGrammarFilter_FilterWithNonFilledGrammar(t *testing.T) {
	filter := NewFilledGrammarFilter()

	grammar, err := filter.Filter(inlinedgrammary.New(
		grammary.NewSymbol("<F>"),
		inlinedgrammary.NewRule(grammary.NewSymbol("<F>"), []grammary.Symbol{
			grammary.NewSymbol("<S>"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<S>"), []grammary.Symbol{
			grammary.NewSymbol("i"),
			grammary.NewSymbol("a"),
			grammary.NewSymbol("<F>"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<E>"), []grammary.Symbol{
			grammary.NewSymbol("<E>"),
			grammary.NewSymbol("+"),
			grammary.NewSymbol("<E>"),
		}),
	))

	assert.NoError(t, err)

	assert.Equal(
		t,
		"{{|F|} [{{|F|} [{<F>}]} {{<F>} [{<S>}]} {{<S>} [{i} {a} {<F>}]} {{<E>} [{<E>} {+} {<E>}]}]}",
		fmt.Sprint(grammar),
	)
}

func TestFilledGrammarFilter_FilterWithFilledGrammar(t *testing.T) {
	filter := NewFilledGrammarFilter()

	grammar, err := filter.Filter(inlinedgrammary.New(
		grammary.NewSymbol("<F>"),
		inlinedgrammary.NewRule(grammary.NewSymbol("<F>"), []grammary.Symbol{
			grammary.NewSymbol("<S>"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<S>"), []grammary.Symbol{
			grammary.NewSymbol("i"),
			grammary.NewSymbol("a"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<E>"), []grammary.Symbol{
			grammary.NewSymbol("<E>"),
			grammary.NewSymbol("+"),
			grammary.NewSymbol("<E>"),
		}),
	))

	assert.NoError(t, err)

	assert.Equal(
		t,
		"{{<F>} [{{<F>} [{<S>}]} {{<S>} [{i} {a}]} {{<E>} [{<E>} {+} {<E>}]}]}",
		fmt.Sprint(grammar),
	)
}
