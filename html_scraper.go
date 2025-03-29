package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func getHTML(rawURL string) (string, error) {
	response, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}

	if response.StatusCode >= 400 {
		errMsg := fmt.Sprintf("%v returned status code: %v", rawURL, response.StatusCode)
		return "", errors.New(errMsg)
	}

	contentType := response.Header.Get("content-type")
	if !strings.Contains(contentType, "text/html") {
		errMsg := fmt.Sprintf("%v returned content-type: %v", rawURL, contentType)
		return "", errors.New(errMsg)
	}

	defer response.Body.Close()

	htmlBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(htmlBody), nil
}

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	var result []string
	doc, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, a := range n.Attr {
				if a.Key == "href" {

					ref, err := url.Parse(a.Val)
					if err != nil {
						return nil, err
					}

					if ref.Path == "/" {
						result = append(result, rawBaseURL)
					} else if ref.Host == "" {
						urlResult, _ := url.JoinPath(rawBaseURL, ref.Path)
						result = append(result, urlResult)
					} else {
						result = append(result, a.Val)
					}

				}
			}
		}
	}
	return result, nil
}
