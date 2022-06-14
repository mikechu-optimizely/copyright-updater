package source

import (
	"bufio"
	"fmt"
	"os"
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

func (r *CodeFile) readLines() error {
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
		if tmp := scanner.Text(); len(tmp) != 0 {
			r.Contents = append(r.Contents, tmp)
		}
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
