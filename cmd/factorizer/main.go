package main

import (
	"compiler/pkg/common/grammar"
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

	grammatic, err := grammar.Parse(rawData)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(grammar.Serialize(grammatic))
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
