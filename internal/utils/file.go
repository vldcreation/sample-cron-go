package utils

import (
	"errors"
	"os"
	"strings"
)

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func ParseFile(filename string) (string, string, error) {
	// parse filename, e.g. "file.txt" => "file", "txt"
	_ext := strings.LastIndex(filename, ".")
	_name := strings.LastIndex(filename, "/") // assume filename inside a folder

	if _ext == -1 {
		return "", "", errors.New("invalid filename")
	}

	name := filename[_name+1 : _ext]
	ext := filename[_ext+1:]

	return name, ext, nil
}
