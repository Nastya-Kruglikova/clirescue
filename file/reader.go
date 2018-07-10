package file

import (
	"errors"
	"io/ioutil"
)

var ErrorInvalidPath = errors.New("Invalid path to file")

func Read(path string) ([]byte, error) {
	if path == "" {
		return nil, ErrorInvalidPath
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}
