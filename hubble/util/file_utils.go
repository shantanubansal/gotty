package util

import (
	"io/ioutil"
	"os"
)

func IsFileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	} else {
		return true, nil
	}
}

func ReadFiles(dirPath string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(dirPath)
	return files, err
}

func DeleteFile(path string) error {
	if isExist, err := IsFileExists(path); err != nil {
		return err
	} else if isExist {
		if err = os.Remove(path); err != nil {
			return err
		}
	}
	return nil
}
