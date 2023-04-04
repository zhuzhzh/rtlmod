package helper

import (
	"io/ioutil"
	"os"
	"strings"
)

func CreateOutputDir(outDir string) error {
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err = os.MkdirAll(outDir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadFiles(fileList string) ([]string, error) {
	filesData, err := ioutil.ReadFile(fileList)
	if err != nil {
		return nil, err
	}

	files := strings.Split(string(filesData), "\n")
	return files, nil
}
