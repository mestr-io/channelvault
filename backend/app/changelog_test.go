package app

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

const sampleHTML = `
<html>
  <body>
    <section id="final-cta">
      <div>Test Content</div>
    </section>
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

	section := FindSectionByID(doc, "final-cta")
	if section == nil {
		t.Fatal("Section with id 'final-cta' not found")
	}
	if section.Data != "section" {
		t.Errorf("Expected 'section', got '%s'", section.Data)
	}

	// Check that the correct section is found
	found := false
	for _, attr := range section.Attr {
		if attr.Key == "id" && attr.Val == "final-cta" {
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

	section := FindSectionByID(doc, "final-cta")
	if section == nil {
		t.Fatal("Section with id 'final-cta' not found")
	}

	htmlStr := RenderNode(section)
	if !strings.Contains(htmlStr, "Test Content") {
		t.Errorf("Rendered HTML does not contain expected content: %s", htmlStr)
	}
	if !strings.HasPrefix(htmlStr, "<section") {
		t.Errorf("Rendered HTML does not start with <section: %s", htmlStr)
	}
}
