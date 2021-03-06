package main

import (
	"context"
	"io"
	stdlog "log"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	dataslr "compiler/data/SLR"
	"compiler/pkg/ast/infrastructure/backend"
	"compiler/pkg/ast/infrastructure/builder"
	"compiler/pkg/ast/infrastructure/executor"
	commonlexer "compiler/pkg/common/lexer"
	"compiler/pkg/lexer/infrastructure"
	lexerexecutor "compiler/pkg/lexer/infrastructure/executor"
	"compiler/pkg/slr"
	"compiler/pkg/slr/export"
	filter "compiler/pkg/slr/filter/app"
	slrgenerator "compiler/pkg/slr/generator/app"
	parserapp "compiler/pkg/slr/parser/app"
	slrrunnner "compiler/pkg/slr/runner/app"
	serializer "compiler/pkg/slr/serializer/app"
)

func main() {
	ctx := context.Background()

	err := runApp(ctx, os.Args)
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			stdlog.Println(string(exitErr.Stderr))
		}
		stdlog.Fatal(err)
	}
}

func runApp(ctx context.Context, args []string) error {
	app := cli.App{
		Name:   "compiler",
		Action: executeAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "lexer",
				Usage:   "Lexer executable path",
				EnvVars: []string{"LEXER_EXECUTABLE"},
			},
			&cli.BoolFlag{
				Name: "printAST",
				EnvVars: []string{
					"PRINT_AST",
				},
			},
		},
	}

	return app.RunContext(ctx, args)
}

func executeAction(ctx *cli.Context) error {
	parser := parserapp.NewInlinedGrammarParser()

	grammar, err := parser.Parse(dataslr.EmbeddedGrammar())
	if err != nil {
		return err
	}

	lexer, err := initLexer(ctx, os.Stdin)
	if err != nil {
		return err
	}

	analyzer := buildAnalyzer(lexer)

	stack, err := analyzer.Analyze(grammar)

	astBackendExecutor, err := executor.NewASTBackendExecutor()
	if err != nil {
		return err
	}

	jsBuilder := builder.NewJsBuilder(backend.NewJSASTBuilderBackend(astBackendExecutor), ctx.Bool("printAST"))

	program, err := jsBuilder.Build(stack)
	if err != nil {
		return err
	}

	_, _ = os.Stdout.WriteString(string(program))

	return err
}

func buildAnalyzer(lexer commonlexer.Lexer) slr.Analyzer {
	filledGrammarFilter := filter.NewFilledGrammarFilter()
	emptySymbolFilter := filter.NewEmptySymbolFilter()

	generator := slrgenerator.NewGenerator()
	validator := slrgenerator.NewValidator()
	tableSerializer := serializer.NewCSVSerializer()
	exporter := export.NewOptionalExporter(
		export.NewFileTableExporter(tableSerializer, "data/SLR/table.csv"),
		false,
	)
	runner := slrrunnner.NewRunner(lexer)

	return slr.NewAnalyzer(
		[]filter.InlinedGrammarFilter{
			filledGrammarFilter,
			emptySymbolFilter,
		},
		generator,
		validator,
		exporter,
		runner,
	)
}

func initLexer(ctx *cli.Context, programReader io.Reader) (infrastructure.LexerAdapter, error) {
	lexerExecutable := ctx.String("lexer")
	if lexerExecutable == "" {
		return nil, errors.New("no lexerExecutable passed")
	}

	executor := lexerexecutor.NewLexerExecutor(lexerExecutable)

	err := executor.Start()
	if err != nil {
		return nil, err
	}

	rawProgram, err := io.ReadAll(programReader)
	if err != nil {
		return nil, err
	}

	if len(rawProgram) == 0 {
		return nil, errors.New("no program or empty file passed")
	}

	err = executor.Write(string(rawProgram))
	if err != nil {
		return nil, err
	}

	return infrastructure.NewLexerAdapter(executor), nil
}
