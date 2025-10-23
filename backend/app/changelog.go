package app

import (
	"io"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// URL to fetch the changelog from
const url = "https://mestr.io/changelog.html"
const sectionName = "div"
const sectionId = "changelog-entries"

// FetchPage fetches the HTML content from the specified URL.
func (a *App) FetchPage() (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// FetchChangelog fetches the changelog and returns the HTML of the section with id="final-cta".
func (a *App) FetchChangelog() (string, error) {
	body, err := a.FetchPage()
	if err != nil {
		return "", err
	}

	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return "", err
	}

	section := FindSectionByID(doc, sectionId)
	if section == nil {
		return "", nil // or return an error if you prefer
	}

	html := RenderNode(section)

	return html, nil
}

// FindSectionByID traverses the HTML node tree to find a section with the specified id.
func FindSectionByID(node *html.Node, id string) *html.Node {
	if node.Type == html.ElementNode && node.Data == sectionName {
		for _, attr := range node.Attr {
			if attr.Key == "id" && attr.Val == id {
				return node
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if result := FindSectionByID(child, id); result != nil {
			return result
		}
	}
	return nil
}

// RenderNode renders an clean HTML node and its subtree to a string.
func RenderNode(node *html.Node) string {
	var builder strings.Builder
	html.Render(&builder, node)

	html := builder.String()

	// Add that to trim all newlines and extra spaces
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(html, " ")
	return cleaned
}
