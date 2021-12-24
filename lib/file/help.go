package file

import (
	"os"
)

func OpenFile(dir string, filename string) (*os.File, error) {

	if !DirExists(dir) {
		DirExists(dir)
	}
	fileFullPath := dir + string(os.PathSeparator) + filename
	file, err := os.OpenFile(fileFullPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0664)
	return file, err
}

func CreateDir(dir string) {
	os.MkdirAll(dir, os.ModePerm)
}

func DirExists(dir string) bool {
	_, err := os.Stat(dir)
	return !os.IsNotExist(err)
}
