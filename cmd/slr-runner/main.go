package main

import (
	"compiler/pkg/common/grammary"
	"compiler/pkg/slr"
	"compiler/pkg/slr/common/inlinedgrammary"
	"compiler/pkg/slr/export"
	slrgenerator "compiler/pkg/slr/generator/app"
	serializer "compiler/pkg/slr/serializer/app"
	"github.com/urfave/cli/v2"
	stdlog "log"
	"os"
)

func main() {
	err := runApp(os.Args)
	if err != nil {
		stdlog.Fatal(err)
	}
}

func runApp(args []string) error {
	app := cli.App{
		Name: "slr-runner",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "lexer",
				Aliases: []string{"l"},
			},
			&cli.StringFlag{
				Name:    "input",
				Aliases: []string{"i"},
			},
			&cli.StringFlag{
				Name:    "grammar",
				Aliases: []string{"g"},
			},
		},
		EnableBashCompletion: true,
		Action:               executeAction,
	}

	return app.Run(args)
}

func executeAction(ctx *cli.Context) error {
	analyzer := buildAnalyzer(ctx)

	//table, err := generator.GenerateTable(inlinedgrammary.New(
	//	grammary.NewSymbol("<F>"),
	//	inlinedgrammary.NewRule(grammary.NewSymbol("<F>"), []grammary.Symbol{
	//		grammary.NewSymbol("<S>"),
	//	}),
	//	inlinedgrammary.NewRule(grammary.NewSymbol("<S>"), []grammary.Symbol{
	//		grammary.NewSymbol("("),
	//		grammary.NewSymbol("<S>"),
	//		grammary.NewSymbol(")"),
	//	}),
	//	inlinedgrammary.NewRule(grammary.NewSymbol("<S>"), []grammary.Symbol{
	//		grammary.NewSymbol("("),
	//		grammary.NewSymbol(")"),
	//	}),
	//))

	// Simple grammar
	//table, err := generator.GenerateTable(inlinedgrammary.New(
	//	grammary.NewSymbol("<F>"),
	//	inlinedgrammary.NewRule(grammary.NewSymbol("<F>"), []grammary.Symbol{
	//		grammary.NewSymbol("a"),
	//	}),
	//))

	err := analyzer.Analyze(inlinedgrammary.New(
		grammary.NewSymbol("<F>"),
		inlinedgrammary.NewRule(grammary.NewSymbol("<F>"), []grammary.Symbol{
			grammary.NewSymbol("<S>"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<S>"), []grammary.Symbol{
			grammary.NewSymbol("a"),
			grammary.NewSymbol("b"),
		}),
		inlinedgrammary.NewRule(grammary.NewSymbol("<S>"), []grammary.Symbol{
			grammary.NewSymbol("c"),
		}),
	))

	return err
}

func buildAnalyzer(_ *cli.Context) slr.Analyzer {
	generator := slrgenerator.NewGenerator()
	validator := slrgenerator.NewValidator()
	tableSerializer := serializer.NewCSVSerializer()
	exporter := export.NewOptionalExporter(
		export.NewFileTableExporter(tableSerializer, "table.csv"),
		true,
	)

	return slr.NewAnalyzer(
		generator,
		validator,
		exporter,
	)
}
