package main

import (
	"compiler/pkg/common/lexer"
	"compiler/pkg/lexer/infrastructure"
	lexerexecutor "compiler/pkg/lexer/infrastructure/executor"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	config, err := parseConfig()
	if err != nil {
		log.Fatalf("Parse config failed %s", err.Error())
	}

	rawData, err := getFileData(config.InputProgramFilePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	lexerImpl, err := initLexer(config, rawData)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer lexerImpl.Close()

	err = runModule(lexerImpl)
	if err != nil {
		log.Fatal(err)
	}
}

func runModule(lexerImpl lexer.Lexer) error {
	for {
		lexem, err := lexerImpl.FetchLexem()
		if err != nil {
			return err
		}
		if lexem.Type == lexer.LexemTypeEOF {
			fmt.Println("program ends")
			break
		}

		fmt.Println("LEXEM:", lexem)
	}

	return nil
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
