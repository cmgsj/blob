package main

import (
	"fmt"
	"os"

	"github.com/cmgsj/blob/pkg/cmd/blob"
)

func main() {
	err := blob.NewCmdBlob().Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
