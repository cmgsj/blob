package storage

const (
	ErrBlobNotFound = errorString("blob not found")
)

type errorString string

func (e errorString) Error() string {
	return string(e)
}
