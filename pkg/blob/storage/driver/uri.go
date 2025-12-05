package driver

import (
	"fmt"
	"net/url"
	"strings"
)

type URI struct {
	Host         string
	Bucket       string
	ObjectPrefix string
}

func ParseURI(driverType, rawURI string) (*URI, error) {
	u, err := url.Parse(rawURI)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" {
		return nil, fmt.Errorf("invalid %s uri %q: scheme is required", driverType, rawURI)
	}

	if u.Host == "" {
		return nil, fmt.Errorf("invalid %s uri %q: host is required", driverType, rawURI)
	}

	path := strings.Split(strings.Trim(u.Path, "/"), "/")

	if len(path) == 0 || path[0] == "" {
		return nil, fmt.Errorf("invalid %s uri %q: path is required", driverType, rawURI)
	}

	uri := &URI{
		Host:   u.Host,
		Bucket: path[0],
	}

	if len(path) > 1 {
		uri.ObjectPrefix = strings.Join(path[1:], "/")
	}

	return uri, nil
}
