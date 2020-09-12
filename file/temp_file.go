package file

import (
	"io/ioutil"
	"os"
)

func WriteTempFile(prefix string, data []byte) (string, error) {
	file, err := ioutil.TempFile(os.TempDir(), prefix)
	if err != nil {
		return "", err
	}
	_, err = file.Write(data)
	if err1 := file.Close(); err == nil {
		err = err1
	}
	if err != nil {
		_ = os.Remove(file.Name())
		return "", err
	}
	return file.Name(), nil
}
