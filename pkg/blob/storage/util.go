package storage

import "strings"

func cleanBlobPrefix(path string) string {
	return strings.Trim(path, "/")
}

func joinBlobPrefix(elems ...string) string {
	var path []string

	for _, elem := range elems {
		elem = strings.Trim(elem, "/")

		if elem != "" {
			path = append(path, elem)
		}
	}
	return strings.Join(path, "/")
}
