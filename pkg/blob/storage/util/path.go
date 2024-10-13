package util

import "strings"

const BlobsPrefix = "blobs"

func BlobObjectPrefix(prefix string) string {
	if prefix == "" {
		return BlobsPrefix
	}

	return BlobPath(prefix, BlobsPrefix)
}

func BlobPath(base string, elems ...string) string {
	var path []string

	base = strings.Trim(base, "/")

	if base != "" {
		path = append(path, base)
	}

	for _, elem := range elems {
		elem = strings.Trim(elem, "/")

		if elem != "" {
			path = append(path, elem)
		}
	}

	return strings.Join(path, "/")
}
