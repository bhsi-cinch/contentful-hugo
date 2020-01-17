package write

import (
	"fmt"
	"os"
	"strings"
)

type Writer struct {
	Store Store
	Files []string
}

func (w *Writer) SaveToFile(fileName string, output string) error {
	var fileMode os.FileMode
	fileMode = 0733

	err := w.Store.MkdirAll(dirForFile(fileName), fileMode)
	if err != nil {
		return err
	}

	bytes := []byte(output)
	err = w.Store.WriteFile(fileName, bytes, fileMode)
	if err != nil {
		return err
	}
	w.Files = append(w.Files, fmt.Sprintf("%s: %d", fileName, len(bytes)))

	return nil
}

func dirForFile(filename string) string {
	index := strings.LastIndex(filename, "/")
	return filename[0:index]
}
