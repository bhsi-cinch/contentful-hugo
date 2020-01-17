package read

type Getter interface {
	Get(url string) (result []byte, err error)
}
