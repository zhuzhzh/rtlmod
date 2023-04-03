package helper

import "os"

func CreateOutputDir(outDir string) error {
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err = os.MkdirAll(outDir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
