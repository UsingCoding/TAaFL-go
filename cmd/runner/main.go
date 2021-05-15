package main

import (
	"compiler/pkg/common/grammary"
	"compiler/pkg/common/lexer"
	"compiler/pkg/common/reporter"
	factorizerapp "compiler/pkg/factorizer/app"
	appgenerator "compiler/pkg/generator/app"
	"compiler/pkg/lexer/infrastructure"
	lexerexecutor "compiler/pkg/lexer/infrastructure/executor"
	"compiler/pkg/runner/app"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
)

func main() {
	reporterImpl := initReporter()

	config, err := parseConfig()
	if err != nil {
		log.Fatalf("Parse config failed %s", err.Error())
	}

	grammar, err := initGrammar(config.GrammarFilePath)
	if err != nil {
		log.Fatalf("Parse grammar failed %s", err.Error())
	}

	program, err := getFileData(config.InputProgramFilePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	lexerImpl, err := initLexer(config, program)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer lexerImpl.Close()

	err = runModule(grammar, lexerImpl, reporterImpl)
	if err != nil {
		log.Fatal(err)
	}
}

func runModule(grammar grammary.Grammar, lexerImpl lexer.Lexer, reporter reporter.Reporter) error {
	grammarWithHeadSequences, err := factorizerapp.BuildHeadSequencesForGrammar(grammar)
	if err != nil {
		return err
	}

	serializer := grammary.NewSerializer()

	serializedGrammar, err := serializer.SerializeGrammar(&grammarWithHeadSequences)
	if err != nil {
		return err
	}

	leftParts, rightParts := appgenerator.CreateTables(serializedGrammar)

	err = app.Runner(leftParts, rightParts, lexerImpl)
	if err != nil && !errors.Is(err, app.AnalyzerErr{}) {
		return err
	}

	msg := "Program text is %s\n"

	if err == nil {
		reporter.Info(fmt.Sprintf(msg, "matched"))
	} else {
		reporter.WithError(err, fmt.Sprintf(msg, "not matched"))
	}

	return nil
}

func initGrammar(path string) (grammary.Grammar, error) {
	data, err := getFileData(path)
	if err != nil {
		return grammary.Grammar{}, err
	}
	return grammary.Parse(data)
}

func initLexer(config *config, rawProgram string) (infrastructure.LexerAdapter, error) {
	executor := lexerexecutor.NewLexerExecutor(config.LexerExecutablePath)

	err := executor.Start()
	if err != nil {
		return nil, err
	}

	err = executor.Write(rawProgram)
	if err != nil {
		return nil, err
	}

	return infrastructure.NewLexerAdapter(executor), nil
}

func getFileData(path string) (string, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func initReporter() reporter.Reporter {
	impl := logrus.New()
	impl.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})

	return reporter.New(impl)
}
