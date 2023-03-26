package logy

import "os"

type fileReader interface {
	ReadFile(name string) ([]byte, error)
}

type configFileReader struct {
}

func newConfigFileReader() *configFileReader {
	return &configFileReader{}
}

func (r *configFileReader) ReadFile(name string) ([]byte, error) {
	file, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	return file, nil
}
