package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"source"
	"strings"
)

const configJsonPath = "src/config.json"

type ExtensionAndBlock struct {
	FileExtension          string
	CommentBlockPath       string
	CommentBlockText       string
	CommentBlockStartsWith string
	CommentBlockEndsWith   string
}

type Configuration struct {
	RootDirectoryToSearch    string
	ExtensionsAndBlocks      []ExtensionAndBlock
	ExtensionCommentBlockMap map[string]ExtensionAndBlock
}

//var wg = sync.WaitGroup{}
var config Configuration

func main() {
	readErr := readConfig()
	if readErr != nil {
		return
	}

	processExtensionBlocks()

	updateErr := updateFiles(config.RootDirectoryToSearch)
	if updateErr != nil {
		fmt.Printf("Failed to update files with extensions %s in dir %s", config.ExtensionsAndBlocks, config.RootDirectoryToSearch)
		return
	}

	//wg.Wait()
}

func readConfig() error {
	file, err := os.Open(configJsonPath)

	if err != nil {
		panic(fmt.Sprintf("Failure opening config file %s", configJsonPath))
	}

	defer func() {
		closeError := file.Close()
		if err == nil {
			err = closeError
		}
	}()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&config)
	if err != nil {
		panic(fmt.Sprintf("Error while decoding config file: %s", err))
	}

	return nil
}

func processExtensionBlocks() {
	for _, extensionAndBlock := range config.ExtensionsAndBlocks {
		// normalize the malformed extensions from config
		if !strings.HasPrefix(extensionAndBlock.FileExtension, ".") {
			extensionAndBlock.FileExtension = "." + extensionAndBlock.FileExtension
		}

		// read the comment block file
		fileBytes, err := os.ReadFile(extensionAndBlock.CommentBlockPath)
		if err != nil {
			panic(fmt.Sprintf("Unable to read comment block from %s", extensionAndBlock.CommentBlockPath))
		}
		extensionAndBlock.CommentBlockText = strings.TrimSuffix(string(fileBytes), "\n")

		if config.ExtensionCommentBlockMap == nil {
			config.ExtensionCommentBlockMap = make(map[string]ExtensionAndBlock)
		}
		config.ExtensionCommentBlockMap[extensionAndBlock.FileExtension] = extensionAndBlock
	}
}

func updateFiles(rootPath string) error {
	err := filepath.Walk(
		rootPath,
		func(path string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if fileInfo.IsDir() {
				return nil
			}

			configBlock := config.ExtensionCommentBlockMap[filepath.Ext(fileInfo.Name())]
			if configBlock.CommentBlockText == "" {
				return nil
			}

			file, err := os.OpenFile(path, os.O_RDWR, 0600)
			if err != nil {
				fmt.Printf("Unable to open %s\n", path)
				return err
			}

			defer func() {
				closeError := file.Close()
				if err == nil {
					err = closeError
				}
			}()

			//wg.Add(1)
			/*go*/
			updateFile(file, configBlock)

			return nil
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func updateFile(fileToUpdate *os.File, blockConfig ExtensionAndBlock) {
	codeFile := source.NewCodeFile(fileToUpdate.Name())

	err := codeFile.RemoveFirstCommentBlock(blockConfig.CommentBlockStartsWith, blockConfig.CommentBlockEndsWith)
	if err != nil {
		panic(fmt.Sprintf("Error while removing disclaimer to %s\n%s", fileToUpdate.Name(), err.Error()))
	}

	err = codeFile.Prepend(blockConfig.CommentBlockText)
	if err != nil {
		panic(fmt.Sprintf("Error while adding disclaimer to %s\n%s", fileToUpdate.Name(), err.Error()))
	}
	fmt.Printf("Added comment block to %s\n", fileToUpdate.Name())

	//wg.Done()
}
