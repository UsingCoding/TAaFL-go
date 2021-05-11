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

	app.RemoveLeftRecursion(grammar)

	grammarSerializer := grammary.NewSerializer()

	serializedGrammar, err := grammarSerializer.SerializeGrammar(&grammar)
	if err != nil {
		return err
	}
	fmt.Println(serializedGrammar)

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
