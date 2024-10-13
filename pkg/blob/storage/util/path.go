package util

import "strings"

func BlobNamePrefix(path string) string {
	return strings.Trim(path, "/")
}

func JoinBlobPath(elems ...string) string {
	var path []string

	for _, elem := range elems {
		elem = strings.Trim(elem, "/")

		if elem != "" {
			path = append(path, elem)
		}
	}

	return strings.Join(path, "/")
}
