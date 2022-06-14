package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const configJsonPath = "src/config.json"

type ExtensionAndBlock struct {
	FileExtension    string
	CommentBlockPath string
	CommentBlockText string
}

type Configuration struct {
	RootDirectoryToSearch    string
	ExtensionsAndBlocks      []ExtensionAndBlock
	ExtensionCommentBlockMap map[string]string
}

var wg = sync.WaitGroup{}
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

	wg.Wait()
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
		extension := extensionAndBlock.FileExtension
		// normalize the malformed extensions from config
		if !strings.HasPrefix(extension, ".") {
			extension = "." + extension
		}

		// read the comment block file
		fileBytes, err := os.ReadFile(extensionAndBlock.CommentBlockPath)
		if err != nil {
			panic(fmt.Sprintf("Unable to read comment block from %s", extensionAndBlock.CommentBlockPath))
		}

		if config.ExtensionCommentBlockMap == nil {
			config.ExtensionCommentBlockMap = make(map[string]string)
		}
		config.ExtensionCommentBlockMap[extension] = string(fileBytes)
	}
}

func updateFiles(rootPath string) error {
	err := filepath.WalkDir(
		rootPath,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			commentBlock := config.ExtensionCommentBlockMap[filepath.Ext(d.Name())]
			if commentBlock == "" {
				return nil
			}

			filePath := rootPath + d.Name()
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

			wg.Add(1)
			go updateFile(file, commentBlock)

			return nil
		},
	)
	if err != nil {
		panic("Error while updating files")
	}

	return nil
}

func updateFile(fileToUpdate *os.File, commentBlock string) {
	tryRemoveExistingDisclaimer(fileToUpdate)
	addDisclaimer(fileToUpdate, commentBlock)
	wg.Done()
}

func tryRemoveExistingDisclaimer(fileToUpdate *os.File) {
	time.Sleep(1 * time.Second) // simulating
	fmt.Printf("Removed disclaimer from %s\n", fileToUpdate.Name())
}

func addDisclaimer(fileToUpdate *os.File, commentBlock string) {
	time.Sleep(2 * time.Second) // simulating
	fmt.Printf("Added disclaimer to %s\n", fileToUpdate.Name())
	fmt.Println(commentBlock)
}
