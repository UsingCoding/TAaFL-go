package main

import (
	"compiler/pkg/common/grammary"
	"compiler/pkg/factorizer/app"
	"fmt"
	"io/ioutil"
)

const (
	exampleFilePath = "data/factorizer/grammar-example"
)

func main() {
	rawData, err := getFileData(exampleFilePath)
	if err != nil {
		fmt.Println(err)
	}

	grammar, err := grammary.Parse(rawData)
	if err != nil {
		fmt.Println(err)
	}

	grammarSerializer := grammary.NewSerializer()

	grammarWithHeadSequences, err := app.BuildHeadSequencesForGrammar(grammar)
	if err != nil {
		fmt.Println(err)
	}

	serializedGrammar, err := grammarSerializer.SerializeGrammar(&grammarWithHeadSequences)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(serializedGrammar)
}

func getFileData(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	// convert bytes to string
	str := string(b)

	// show file data
	return str, nil
}
