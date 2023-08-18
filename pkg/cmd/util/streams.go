package util

import (
	"io"
	"os"
)

type IOStreams struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

func NewOSStreams() IOStreams {
	return IOStreams{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
}
