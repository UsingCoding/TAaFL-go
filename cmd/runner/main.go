package main

import (
	"compiler/pkg/common/grammary"
	"compiler/pkg/common/lexer"
	appgenerator "compiler/pkg/generator/app"
	"compiler/pkg/lexer/infrastructure"
	lexerexecutor "compiler/pkg/lexer/infrastructure/executor"
	"compiler/pkg/runner/app"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
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

	err = runModule(grammar, lexerImpl)
	if err != nil {
		log.Fatal(err)
	}
}

func runModule(grammar grammary.Grammar, lexerImpl lexer.Lexer) error {
	serializer := grammary.NewSerializer()

	serializedGrammar, err := serializer.SerializeGrammar(grammar)
	if err != nil {
		return err
	}

	leftParts, rightParts := appgenerator.CreateTables(serializedGrammar)

	isMatched, err := app.Runner(leftParts, rightParts, lexerImpl)
	if err != nil {
		return err
	}

	msg := "Program text is %s\n"

	if isMatched {
		fmt.Printf(msg, "matched")
	} else {
		fmt.Printf(msg, "not matched")
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
