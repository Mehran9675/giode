package markdown

import (
	"image"
	"strings"
	"testing"

	"gioui.org/layout"
	"gioui.org/op"

	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/fonts"
)

const doc = "# Title\n\n" +
	"Run **go test** or `go vet`, read the [docs](https://gioui.org), ~~old~~.\n\n" +
	"- one\n- two\n\n" +
	"1. first\n2. second\n\n" +
	"- [x] done\n- [ ] todo\n\n" +
	"> quoted\n\n" +
	"```go\nfunc main() {\n\tprintln(\"hi\")\n}\n```\n\n" +
	"| a | b |\n|---|--:|\n| 1 | 2 |\n\n---\n"

func TestMarkdown_RendersBlocks(t *testing.T) {
	fonts.UseGofont()
	elements.SetWindowSize(image.Pt(600, 2000))
	defer elements.SetWindowSize(image.Point{})
	elements.BeginFrame()
	elements.BeginSemantics()
	var ops op.Ops
	dims := New(doc).Layout(layout.Context{Ops: &ops, Constraints: layout.Constraints{Max: image.Pt(600, 1e6)}})
	nodes := elements.EndSemantics()

	if dims.Size.Y <= 0 || dims.Size.Y > 2000 {
		t.Fatalf("document height %d", dims.Size.Y)
	}
	var texts []string
	copyButtons := 0
	for _, n := range nodes {
		texts = append(texts, n.Label)
		if n.Kind == elements.SemanticButton && strings.Contains(n.Label, "Copy") {
			copyButtons++
		}
	}
	all := strings.Join(texts, "|")
	for _, want := range []string{
		"Title",
		"Run go test or go vet, read the docs, old.",
		"one", "second", "done", "todo", "quoted",
		"func main() {\n    println(\"hi\")\n}",
		"go",
	} {
		if !strings.Contains(all, want) {
			t.Errorf("rendered text lacks %q; got %q", want, all)
		}
	}
	if copyButtons != 1 {
		t.Errorf("%d copy buttons, want 1 for the code block", copyButtons)
	}
}

func TestMarkdown_ParseCached(t *testing.T) {
	a, b := parse("hello *world*"), parse("hello *world*")
	if a != b {
		t.Error("same source parsed twice")
	}
}
