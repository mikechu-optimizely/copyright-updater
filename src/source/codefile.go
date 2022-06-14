package source

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CodeFile struct {
	Filename string
	Contents []string
}

func NewCodeFile(filename string) *CodeFile {
	return &CodeFile{
		Filename: filename,
		Contents: make([]string, 0),
	}
}

// Returns true if a comment block was found on the first line of the code
// Function will halt if a block is not found
func (r *CodeFile) readLinesExceptTopCommentBlock(commentBlockStartToken string, commentBlockEndToken string) (bool, error) {
	if _, err := os.Stat(r.Filename); err != nil {
		return false, nil
	}

	f, err := os.OpenFile(r.Filename, os.O_RDONLY, 0600)
	if err != nil {
		return false, err
	}

	defer func() {
		closeError := f.Close()
		if err == nil {
			err = closeError
		}
	}()

	foundFirstCommentBlock := false
	ignoringLines := false
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		tmp := scanner.Text()

		if !strings.Contains(tmp, commentBlockStartToken) && !foundFirstCommentBlock {
			fmt.Println("No comment block found on line 1 of file so skipping removal")
			return false, nil
		}

		if strings.Contains(tmp, commentBlockStartToken) {
			foundFirstCommentBlock = true
			ignoringLines = true
			fmt.Println("Found comment block")
		}

		if !ignoringLines {
			r.Contents = append(r.Contents, tmp)
		}

		if strings.Contains(tmp, commentBlockEndToken) && foundFirstCommentBlock {
			ignoringLines = false
		}
	}

	return true, nil
}

func (r *CodeFile) readLines() error {
	alreadyReadContents := len(r.Contents) > 0
	if alreadyReadContents {
		return nil
	}

	if _, err := os.Stat(r.Filename); err != nil {
		return nil
	}

	f, err := os.OpenFile(r.Filename, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}

	defer func() {
		closeError := f.Close()
		if err == nil {
			err = closeError
		}
	}()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		tmp := scanner.Text()
		r.Contents = append(r.Contents, tmp)
	}

	return nil
}

func (r *CodeFile) Prepend(content string) error {
	err := r.readLines()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(r.Filename, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer func() {
		closeError := f.Close()
		if err == nil {
			err = closeError
		}
	}()

	writer := bufio.NewWriter(f)
	_, err = writer.WriteString(fmt.Sprintf("%s\n", content))
	if err != nil {
		return err
	}

	for _, line := range r.Contents {
		_, err := writer.WriteString(fmt.Sprintf("%s\n", line))
		if err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}

func (r *CodeFile) RemoveFirstCommentBlock(startsWith string, endsWith string) error {
	foundBlock, err := r.readLinesExceptTopCommentBlock(startsWith, endsWith)
	if err != nil {
		return err
	}
	if !foundBlock {
		return nil
	}

	f, err := os.OpenFile(r.Filename, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	defer func() {
		closeError := f.Close()
		if err == nil {
			err = closeError
		}
	}()

	writer := bufio.NewWriter(f)
	for _, line := range r.Contents {
		_, err := writer.WriteString(fmt.Sprintf("%s\n", line))
		if err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}
