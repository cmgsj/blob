package util

import (
	"fmt"
	"io"
	"os"
)

func CheckErr(err error, stderr io.Writer) {
	if err != nil {
		fmt.Fprintln(stderr, err)
		os.Exit(1)
	}
}
