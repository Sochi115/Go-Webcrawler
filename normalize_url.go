package main

import (
	"net/url"
	"strings"
)

func normalizeUrl(requestUrl string) (string, error) {
	u, err := url.Parse(requestUrl)
	if err != nil {
		return "", err
	}

	trimmedPath := strings.TrimSuffix(u.Path, "/")
	result := u.Host + trimmedPath
	return result, err
}
