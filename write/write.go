package write

import (
	"os"
	"strings"
)

type Writer struct {
	Store Store
}

func (w *Writer) SaveToFile(fileName string, output string) error {
	var fileMode os.FileMode
	fileMode = 0733

	err := w.Store.MkdirAll(dirForFile(fileName), fileMode)
	if err != nil {
		return err
	}

	err = w.Store.WriteFile(fileName, []byte(output), fileMode)
	if err != nil {
		return err
	}
	return nil
}

func dirForFile(filename string) string {
	index := strings.LastIndex(filename, "/")
	return filename[0:index]
}
