package read

type Reader struct {
	Store Store
}

func (r *Reader) ViewFromFile(fileName string) (body string, err error) {
	result, err := r.Store.ReadFromFile(fileName)

	return string(result), err
}
