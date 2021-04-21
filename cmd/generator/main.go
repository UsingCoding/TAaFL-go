package main

import (
	"compiler/pkg/common/grammary"
	"compiler/pkg/factorizer/app"
	appgenerator "compiler/pkg/generator/app"
	"fmt"
	"io/ioutil"
)

const (
	exampleFilePath = "data/factorizer/grammar-example"
)

func main() {
	err := runService()
	if err != nil {
		fmt.Println(err)
	}
}

func runService() error {
	rawData, err := getFileData(exampleFilePath)
	if err != nil {
		return err
	}

	grammar, err := grammary.Parse(rawData)
	if err != nil {
		return err
	}

	err = app.FactorizeGrammar(grammar)
	if err != nil {
		return err
	}

	grammarSerializer := grammary.NewSerializer()

	grammarWithHeadSequences, err := app.BuildHeadSequencesForGrammar(grammar)
	if err != nil {
		return err
	}

	serializedGrammar, err := grammarSerializer.SerializeGrammar(&grammarWithHeadSequences)
	if err != nil {
		return err
	}
	fmt.Println(serializedGrammar)

	leftParts, rightParts := appgenerator.CreateTables(serializedGrammar)

	fmt.Println(leftParts, rightParts)

	return nil
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
