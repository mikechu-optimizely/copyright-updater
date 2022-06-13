package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const configJsonPath = "src/config.json"

type Configuration struct {
	FileExtensionsToUpdate                  []string
	RelativePathToDisclaimerCommentBlockTxt string
	RelativeRootDirectoryOfFilesToUpdate    string
}

func main() {
	configuration, readErr := readConfig(configJsonPath)
	if readErr != nil {
		fmt.Printf("Unable to read or decode %s", configJsonPath)
		return
	}

	configuration.FileExtensionsToUpdate = normalizeFileExtensions(configuration.FileExtensionsToUpdate)

	commentBlock, commentBlockError := readCommentBlock(configuration.RelativePathToDisclaimerCommentBlockTxt)
	if commentBlockError != nil {
		fmt.Printf("Unable to read comment block from %s", configuration.RelativePathToDisclaimerCommentBlockTxt)
		return
	}

	updateErr := updateFiles(configuration.RelativeRootDirectoryOfFilesToUpdate, configuration.FileExtensionsToUpdate, commentBlock)
	if updateErr != nil {
		fmt.Printf("Failed to update files with extensions %s in dir %s", configuration.FileExtensionsToUpdate, configuration.RelativeRootDirectoryOfFilesToUpdate)
		return
	}
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

func normalizeFileExtensions(fileExtensions []string) []string {
	for index, extension := range fileExtensions {
		if strings.HasPrefix(extension, ".") {
			continue
		}
		fileExtensions[index] = "." + extension
	}
	fmt.Println(fileExtensions)
	return fileExtensions
}

func readCommentBlock(pathToDisclaimerCommentBlock string) (string, error) {
	fileBytes, err := os.ReadFile(pathToDisclaimerCommentBlock)

	if err != nil {
		return "", err
	}

	return string(fileBytes), err
}

func updateFiles(rootPath string, extensionsToUpdate []string, commentBlock string) error {

	filePath := "src/fakeFile.go"
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Unable to Open() %s\n", filePath)
		return err
	}

	defer func() {
		closeError := file.Close()
		if err == nil {
			err = closeError
		}
	}()

	updateFile(file, commentBlock)

	return nil
}

func updateFile(fileToUpdate *os.File, commentBlock string) {
	tryRemoveExistingDisclaimer(fileToUpdate)
	addDisclaimer(fileToUpdate, commentBlock)
}

func tryRemoveExistingDisclaimer(fileToUpdate *os.File) {
	fmt.Printf("tryRemoveExistingDisclaimer(%s)\n", fileToUpdate.Name())
}

func addDisclaimer(fileToUpdate *os.File, block string) {
	fmt.Printf("addDisclaimer(%s)\n", fileToUpdate.Name())
}
