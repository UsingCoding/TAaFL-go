package slr

import (
	"compiler/pkg/slr/common/inlinedgrammary"
	"compiler/pkg/slr/export"
	generator "compiler/pkg/slr/generator/app"
	runner "compiler/pkg/slr/runner/app"
)

type Analyzer interface {
	Analyze(grammar inlinedgrammary.Grammar) error
}

func NewAnalyzer(
	tableGenerator generator.Generator,
	tableValidator generator.Validator,
	tableExporter export.TableExporter,
	runner runner.Runner,
) Analyzer {
	return &analyzer{
		tableGenerator: tableGenerator,
		tableValidator: tableValidator,
		tableExporter:  tableExporter,
		runner:         runner,
	}
}

type analyzer struct {
	tableGenerator generator.Generator
	tableValidator generator.Validator
	tableExporter  export.TableExporter
	runner         runner.Runner
}

func (slr *analyzer) Analyze(grammar inlinedgrammary.Grammar) error {
	table, err := slr.tableGenerator.GenerateTable(grammar)
	if err != nil {
		return err
	}

	err = slr.tableExporter.Export(table)
	if err != nil {
		return err
	}

	err = slr.tableValidator.Validate(table)
	if err != nil {
		return err
	}

	err = slr.runner.Proceed(table, grammar.Axiom())
	if err != nil {
		return err
	}

	return nil
}
