package main

import (
	"flag"
	"github.com/pkg/errors"
)

type config struct {
	GrammarFilePath      string
	LexerExecutablePath  string
	InputProgramFilePath string
}

func parseConfig() (*config, error) {
	lexerPath := flag.String("l", "", "Lexer executable path")
	inputProgramPath := flag.String("i", "", "Input program path")
	grammarFilePath := flag.String("g", "", "Grammar file path")

	flag.Parse()

	if *lexerPath == "" || *inputProgramPath == "" || *grammarFilePath == "" {
		return nil, errors.New("some of flags is missed")
	}

	return &config{
		LexerExecutablePath:  *lexerPath,
		InputProgramFilePath: *inputProgramPath,
	}, nil
}
