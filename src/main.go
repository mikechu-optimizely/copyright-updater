package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const configJsonPath = "src/config.json"

type Configuration struct {
	FileExtensionsToUpdate                  []string
	RelativePathToDisclaimerCommentBlockTxt string
}

func main() {
	configuration, readErr := readConfig(configJsonPath)
	if readErr != nil {
		fmt.Printf("Unable to read or decode %s", configJsonPath)
		return
	}

	fmt.Println(configuration.RelativePathToDisclaimerCommentBlockTxt)
	commentBlock, commentBlockError := readCommentBlock(configuration.RelativePathToDisclaimerCommentBlockTxt)
	if commentBlockError != nil {
		fmt.Printf("Unable to read comment block from %s", configuration.RelativePathToDisclaimerCommentBlockTxt)
		return
	}
	fmt.Println(commentBlock)

	fmt.Println(configuration.FileExtensionsToUpdate)
}

func readConfig(configJsonPath string) (*Configuration, error) {
	file, err := os.Open(configJsonPath)

	if err != nil {
		return nil, err
	}

	defer func() {
		closeError := file.Close()
		if err == nil {
			err = closeError
		}
	}()

	decoder := json.NewDecoder(file)

	configuration := Configuration{}
	err = decoder.Decode(&configuration)

	if err != nil {
		return nil, err
	}

	return &configuration, err
}

func readCommentBlock(pathToDisclaimerCommentBlock string) (string, error) {
	fileBytes, err := os.ReadFile(pathToDisclaimerCommentBlock)

	if err != nil {
		return "", err
	}

	return string(fileBytes), err
}
