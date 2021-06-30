package myvmware

import (
	"fmt"
	"io"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

const samlInputQuery = `input[name="SAMLResponse"]`

var samlInputSelector = cascadia.MustCompile(samlInputQuery)

func getSAMLToken(body io.ReadCloser) (string, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return "", err
	}

	node := samlInputSelector.MatchFirst(doc)
	if node == nil {
		return "", fmt.Errorf("Could not find node that matches %#v", samlInputQuery)
	}

	for _, attr := range node.Attr {
		if attr.Key == "value" {
			return attr.Val, nil
		}
	}

	return "", fmt.Errorf("Could not find the node's value attribute")
}
