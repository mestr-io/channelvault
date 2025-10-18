package app

import (
	"io"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// URL to fetch the changelog from
const url = "https://mestr.io"
const sectionId = "final-cta"

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

	// Add that to triom all newlines and extra spaces
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(html, " ")

	return cleaned, nil
}

func FindSectionByID(n *html.Node, id string) *html.Node {
	if n.Type == html.ElementNode && n.Data == "section" {
		for _, attr := range n.Attr {
			if attr.Key == "id" && attr.Val == id {
				return n
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := FindSectionByID(c, id); result != nil {
			return result
		}
	}
	return nil
}

// RenderNode renders an HTML node and its subtree to a string.
func RenderNode(n *html.Node) string {
	var b strings.Builder
	html.Render(&b, n)
	return b.String()
}
