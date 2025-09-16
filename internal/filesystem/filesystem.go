package filesystem

import (
	"os"
)

const fileDir = "/app/files/"

func Write(name string, content []byte) error {
	err := os.WriteFile(fileDir+name, content, 0600)
	if err != nil {
		return err
	}

	return nil
}

func Read(name string) ([]byte, error) {
	fileContent, err := os.ReadFile(fileDir + name)
	if err != nil {
		return nil, err
	}

	return fileContent, nil
}
