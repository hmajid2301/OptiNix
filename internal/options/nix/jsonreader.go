package nix

import "os"

type JSONReader struct{}

func NewReader() JSONReader {
	return JSONReader{}
}

func (JSONReader) Read(path string) ([]byte, error) {
	nixJSON, err := os.ReadFile(path)
	return nixJSON, err
}
