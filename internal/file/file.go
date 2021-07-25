package file

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
)

func WithFileReader(
	filename string,
	f func(io.Reader) error,
) error {
	file, err := os.Open(filepath.Clean(filename))
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("failed to close input file '%s'", filename)
		}
	}()
	return f(bufio.NewReader(file))
}

func WithFileWriter(
	filename string,
	f func(io.Writer) error,
) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("failed to close output file: %s", err)
		}
	}()

	bufWriter := bufio.NewWriter(file)
	if err := f(bufWriter); err != nil {
		return err
	}
	return bufWriter.Flush()
}
