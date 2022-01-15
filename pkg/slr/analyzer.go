package slr

import (
	"compiler/pkg/ast/app"
	"compiler/pkg/slr/common/inlinedgrammary"
	"compiler/pkg/slr/export"
	filter "compiler/pkg/slr/filter/app"
	generator "compiler/pkg/slr/generator/app"
	runner "compiler/pkg/slr/runner/app"
)

type Analyzer interface {
	Analyze(grammar inlinedgrammary.Grammar) (app.Stack, error)
}

func NewAnalyzer(
	grammarFilters []filter.InlinedGrammarFilter,
	tableGenerator generator.Generator,
	tableValidator generator.Validator,
	tableExporter export.TableExporter,
	runner runner.Runner,
) Analyzer {
	return &analyzer{
		grammarFilters: grammarFilters,
		tableGenerator: tableGenerator,
		tableValidator: tableValidator,
		tableExporter:  tableExporter,
		runner:         runner,
	}
}

type analyzer struct {
	grammarFilters []filter.InlinedGrammarFilter
	tableGenerator generator.Generator
	tableValidator generator.Validator
	tableExporter  export.TableExporter
	runner         runner.Runner
}

func (slr *analyzer) Analyze(grammar inlinedgrammary.Grammar) (app.Stack, error) {
	var filteredGrammar inlinedgrammary.Grammar
	var err error

	for _, grammarFilter := range slr.grammarFilters {
		filteredGrammar, err = grammarFilter.Filter(grammar)
		if err != nil {
			return nil, err
		}
	}

	table, err := slr.tableGenerator.GenerateTable(filteredGrammar)
	if err != nil {
		return nil, err
	}

	err = slr.tableExporter.Export(table)
	if err != nil {
		return nil, err
	}

	err = slr.tableValidator.Validate(table)
	if err != nil {
		return nil, err
	}

	stack, err := slr.runner.Proceed(table, filteredGrammar.Axiom())
	if err != nil {
		return nil, err
	}

	return stack, nil
}
