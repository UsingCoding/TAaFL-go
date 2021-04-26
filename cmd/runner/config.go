package main

import (
	"flag"
	"github.com/pkg/errors"
)

type config struct {
	LexerExecutablePath  string
	InputProgramFilePath string
}

func parseConfig() (*config, error) {
	lexerPath := flag.String("l", "", "Lexer executable path")
	inputProgrammPath := flag.String("f", "", "Input programm path")

	flag.Parse()

	if *lexerPath == "" || *inputProgrammPath == "" {
		return nil, errors.New("some of flags is missed")
	}

	return &config{
		LexerExecutablePath:  *lexerPath,
		InputProgramFilePath: *inputProgrammPath,
	}, nil
}
