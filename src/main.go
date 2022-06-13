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
		fmt.Println("Unable to read or decode config.json")
		return
	}

	fmt.Println(configuration.FileExtensionsToUpdate)
	fmt.Println(configuration.RelativePathToDisclaimerCommentBlockTxt)
}

func readConfig(configJsonPath string) (*Configuration, error) {
	file, err := os.Open(configJsonPath)

	if err != nil {
		fmt.Println("Failed to Open() config.json")
		return nil, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	configuration := Configuration{}
	err = decoder.Decode(&configuration)

	if err != nil {
		fmt.Println("Failed to Decode() configuration file")
		return nil, err
	}

	return &configuration, err
}
