package gzipHelper

import (
	"compress/gzip"
	"io"
	"os"
)

func Unzip(filename, dstFilename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	gz.Close()
	writer, err := os.Create(dstFilename)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, gz)

	return err
}
