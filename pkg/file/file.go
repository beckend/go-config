// Package file file related operations
package file

import (
	os "os"
	filepath "path/filepath"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type WriteFileOptions struct {
	PathFile      string
	ContentsBytes []byte
}

func WriteFile(options *WriteFileOptions) error {
	err := os.MkdirAll(filepath.Dir(options.PathFile), 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(options.PathFile, options.ContentsBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
