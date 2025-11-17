package storage

const (
	ErrBlobNotFound = stringError("blob not found")
)

type stringError string

func (e stringError) Error() string {
	return string(e)
}
