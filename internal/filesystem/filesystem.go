package filesystem

import (
	"io"
	"os"
)

const fileDir = "/app/files/"

func WriteStream(name string, content io.Reader) error {
	file, err := os.Create(fileDir + name)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, content)
	if err != nil {
		return err
	}

	return nil
}

func Read(name string, w io.Writer) error {
	file, err := os.Open(fileDir + name)
	if err != nil {
		return err
	}
	defer file.Close()

	chunk := make([]byte, 1024)

	for {
		n, err := file.Read(chunk)
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		_, err = w.Write(chunk[:n])
		if err != nil {
			return err
		}
	}

	return nil
}
