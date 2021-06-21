package main

import (
	"compiler/pkg/common/lexer"
	filter "compiler/pkg/slr/filter/app"
	slrrunnner "compiler/pkg/slr/runner/app"
	"fmt"
	"github.com/urfave/cli/v2"

	"compiler/pkg/common/grammary"
	"compiler/pkg/slr"
	"compiler/pkg/slr/common/inlinedgrammary"
	"compiler/pkg/slr/export"
	slrgenerator "compiler/pkg/slr/generator/app"
	serializer "compiler/pkg/slr/serializer/app"
)

func executeAction(ctx *cli.Context) error {
	analyzer := buildAnalyzer(ctx)

	// Recursive grammar
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

	if err == nil {
		fmt.Println("OK")
	}

	return err
}

func buildAnalyzer(_ *cli.Context) slr.Analyzer {
	emptySymbolFilter := filter.NewEmptySymbolFilter()
	generator := slrgenerator.NewGenerator()
	validator := slrgenerator.NewValidator()
	tableSerializer := serializer.NewCSVSerializer()
	exporter := export.NewOptionalExporter(
		export.NewFileTableExporter(tableSerializer, "data/SLR/table.csv"),
		true,
	)
	runner := slrrunnner.NewRunner(buildLexer())

	return slr.NewAnalyzer(
		emptySymbolFilter,
		generator,
		validator,
		exporter,
		runner,
	)
}

func buildLexer() lexer.Lexer {
	return &mockLexer{lexems: []lexer.Lexem{
		{
			Type:  lexer.LexemTypeOpenParenthesis,
			Value: "(",
		},
		{
			Type:  lexer.LexemTypeOpenParenthesis,
			Value: "(",
		},
		{
			Type:  lexer.LexemTypeClosingParenthesis,
			Value: ")",
		},
		{
			Type:  lexer.LexemTypeClosingParenthesis,
			Value: ")",
		},
		{
			Type:  lexer.LexemTypeEOF,
			Value: "_|_",
		},
	}}
}
