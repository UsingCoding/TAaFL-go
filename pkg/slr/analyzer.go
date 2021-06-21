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
	grammarFilter generator.InlinedGrammarFilter,
	tableGenerator generator.Generator,
	tableValidator generator.Validator,
	tableExporter export.TableExporter,
	runner runner.Runner,
) Analyzer {
	return &analyzer{
		grammarFilter:  grammarFilter,
		tableGenerator: tableGenerator,
		tableValidator: tableValidator,
		tableExporter:  tableExporter,
		runner:         runner,
	}
}

type analyzer struct {
	grammarFilter  generator.InlinedGrammarFilter
	tableGenerator generator.Generator
	tableValidator generator.Validator
	tableExporter  export.TableExporter
	runner         runner.Runner
}

func (slr *analyzer) Analyze(grammar inlinedgrammary.Grammar) error {
	filteredGrammar, err := slr.grammarFilter.Filter(grammar)
	if err != nil {
		return err
	}

	table, err := slr.tableGenerator.GenerateTable(filteredGrammar)
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

	err = slr.runner.Proceed(table, filteredGrammar.Axiom())
	if err != nil {
		return err
	}

	return nil
}
