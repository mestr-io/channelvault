package app

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

const sampleHTML = `
<html>
  <body>
    <div id="changelog-entries">
      <div>Test     Content</div>
    </div>
    <section id="other">
      <div>Other Content</div>
    </section>
  </body>
</html>
`

func TestFindSectionByID(t *testing.T) {
	doc, err := html.Parse(strings.NewReader(sampleHTML))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	section := FindSectionByID(doc, "changelog-entries")
	if section == nil {
		t.Fatal("div with id 'changelog-entries' not found")
	}
	if section.Data != "div" {
		t.Errorf("Expected 'div', got '%s'", section.Data)
	}

	// Check that the correct section is found
	found := false
	for _, attr := range section.Attr {
		if attr.Key == "id" && attr.Val == "changelog-entries" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Section does not have the correct id attribute")
	}
}

func TestRenderNode(t *testing.T) {
	doc, err := html.Parse(strings.NewReader(sampleHTML))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	section := FindSectionByID(doc, "changelog-entries")
	if section == nil {
		t.Fatal("div with id 'changelog-entries' not found")
	}

	htmlStr := RenderNode(section)
	if !strings.Contains(htmlStr, "Test Content") {
		t.Errorf("Rendered HTML does not contain expected content: %s", htmlStr)
	}
	if !strings.HasPrefix(htmlStr, "<div") {
		t.Errorf("Rendered HTML does not start with <div: %s", htmlStr)
	}
}
