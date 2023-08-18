package util

import (
	"encoding/json"
	"fmt"
	"io"
)

func PrintJSON(out io.Writer, v interface{}) error {
	buf, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	if string(buf) == "{}" {
		buf = []byte("ok")
	}
	fmt.Fprintf(out, "%s\n", buf)
	return nil
	// enc := json.NewEncoder(out)
	// enc.SetIndent("", "  ")
	// return enc.Encode(v)
}
