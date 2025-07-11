package nodes

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/tylertravisty/go-lexical"
)

func TestMain(m *testing.M) {
	exitCode := run(m)
	os.Exit(exitCode)
}

func run(m *testing.M) int {
	return m.Run()
}

func TestUnmarshal(t *testing.T) {
	lexical.ResetNodes()
	lexical.RegisterNodes(&AutoLinkNode{}, &TextNode{}, &ParagraphNode{})
	tests := []string{
		`{"root":{"children":[{"children":[{"detail":0,"format":0,"mode":"normal","style":"","text":"asdf","type":"text","version":1}],"direction":"ltr","format":"","indent":0,"type":"paragraph","version":1,"textFormat":0,"textStyle":""}],"direction":"ltr","format":"","indent":0,"type":"root","version":1}}`,
		`{"root":{"children":[{"children":[{"detail":0,"format":0,"mode":"normal","style":"","text":"with a link? ","type":"text","version":1},{"children":[{"detail":0,"format":0,"mode":"normal","style":"","text":"www.google.com","type":"text","version":1}],"direction":"ltr","format":"","indent":0,"type":"autolink","version":1,"rel":null,"target":null,"title":null,"url":"https://www.google.com","isUnlinked":false},{"detail":0,"format":0,"mode":"normal","style":"","text":" cool!","type":"text","version":1}],"direction":"ltr","format":"","indent":0,"type":"paragraph","version":1,"textFormat":0,"textStyle":""}],"direction":"ltr","format":"","indent":0,"type":"root","version":1}}`,
	}

	for _, test := range tests {
		var root RootNode
		err := json.Unmarshal([]byte(test), &root)
		if err != nil {
			t.Fatal("json.Unmarshal err:", err)
		}

		err = root.Root.Valid()
		if err != nil {
			t.Fatal("root.Root.Valid err:", err)
		}
	}
}

func TestUnmarshalReturnsError(t *testing.T) {
	t.Run("WithUnregisteredNodes", withUnregisteredNodes)
}

func withUnregisteredNodes(t *testing.T) {
	message := `{"root":{"children":[{"children":[{"detail":0,"format":0,"mode":"normal","style":"","text":"with a link? ","type":"text","version":1},{"children":[{"detail":0,"format":0,"mode":"normal","style":"","text":"www.google.com","type":"text","version":1}],"direction":"ltr","format":"","indent":0,"type":"autolink","version":1,"rel":null,"target":null,"title":null,"url":"https://www.google.com","isUnlinked":false},{"detail":0,"format":0,"mode":"normal","style":"","text":" cool!","type":"text","version":1}],"direction":"ltr","format":"","indent":0,"type":"paragraph","version":1,"textFormat":0,"textStyle":""}],"direction":"ltr","format":"","indent":0,"type":"root","version":1}}`

	tests := [][]lexical.Node{
		{},
		{&AutoLinkNode{}},
		{&ParagraphNode{}},
		{&TextNode{}},
		{&ParagraphNode{}, &TextNode{}},
		{&AutoLinkNode{}, &TextNode{}},
		{&AutoLinkNode{}, &ParagraphNode{}},
	}

	for _, test := range tests {
		lexical.ResetNodes()
		lexical.RegisterNodes(test...)
		var root RootNode
		err := json.Unmarshal([]byte(message), &root)
		if err == nil {
			t.Fatal("json.Unmarshal err is nil; expected non-nil err")
		}
	}
}

func TestValidReturnsError(t *testing.T) {
	t.Run("WithInvalidElementNode", withInvalidElementNode)
}

func withInvalidElementNode(t *testing.T) {
	lexical.ResetNodes()
	lexical.RegisterNodes(&ParagraphNode{}, &TextNode{})

	tests := []string{
		`{"root":{"children":[{"children":[{"detail":0,"format":0,"mode":"normal","style":"","text":"asdf","type":"text","version":1}],"direction":"lrt","format":"","indent":0,"type":"paragraph","version":1,"textFormat":0,"textStyle":""}],"direction":"ltr","format":"","indent":0,"type":"root","version":1}}`,
		`{"root":{"children":[{"children":[{"detail":0,"format":0,"mode":"normal","style":"","text":"asdf","type":"text","version":1}],"direction":"rtl","format":"invalid","indent":0,"type":"paragraph","version":1,"textFormat":0,"textStyle":""}],"direction":"ltr","format":"","indent":0,"type":"root","version":1}}`,
		`{"root":{"children":[{"children":[{"detail":0,"format":0,"mode":"normal","style":"","text":"asdf","type":"text","version":1}],"direction":"ltr","format":"","indent":0,"type":"paragraph","version":1,"textFormat":0,"textStyle":""}],"direction":"lrt","format":"","indent":0,"type":"root","version":1}}`,
		`{"root":{"children":[{"children":[{"detail":0,"format":0,"mode":"normal","style":"","text":"asdf","type":"text","version":1}],"direction":"ltr","format":"","indent":0,"type":"paragraph","version":1,"textFormat":0,"textStyle":""}],"direction":"lrt","format":"","indent":0,"type":"root","version":1}}`,
	}

	for _, test := range tests {
		var root RootNode
		err := json.Unmarshal([]byte(test), &root)
		if err != nil {
			t.Fatal("json.Unmarshal err:", err)
		}

		err = root.Root.Valid()
		if err == nil {
			t.Fatal("json.Unmarshal err is nil; expected non-nil err")
		}
	}
}

func TestTextContentSize(t *testing.T) {
	lexical.ResetNodes()
	lexical.RegisterNodes(&AutoLinkNode{}, &ParagraphNode{}, &TextNode{})
	tests := []struct {
		expected int
		message  string
	}{
		{
			expected: 4,
			message:  `{"root":{"children":[{"children":[{"detail":0,"format":0,"mode":"normal","style":"","text":"asdf","type":"text","version":1}],"direction":"ltr","format":"","indent":0,"type":"paragraph","version":1,"textFormat":0,"textStyle":""}],"direction":"ltr","format":"","indent":0,"type":"root","version":1}}`,
		},
		{
			expected: 33,
			message:  `{"root":{"children":[{"children":[{"detail":0,"format":0,"mode":"normal","style":"","text":"with a link? ","type":"text","version":1},{"children":[{"detail":0,"format":0,"mode":"normal","style":"","text":"www.google.com","type":"text","version":1}],"direction":"ltr","format":"","indent":0,"type":"autolink","version":1,"rel":null,"target":null,"title":null,"url":"https://www.google.com","isUnlinked":false},{"detail":0,"format":0,"mode":"normal","style":"","text":" cool!","type":"text","version":1}],"direction":"ltr","format":"","indent":0,"type":"paragraph","version":1,"textFormat":0,"textStyle":""}],"direction":"ltr","format":"","indent":0,"type":"root","version":1}}`,
		},
	}

	for _, test := range tests {
		var root RootNode
		err := json.Unmarshal([]byte(test.message), &root)
		if err != nil {
			t.Fatal("json.Unmarshal err:", err)
		}

		size := root.Root.TextContentSize()
		if size != test.expected {
			t.Fatalf("expected text content size %d; got %d", test.expected, size)
		}
	}

}

func TestNodeTypes(t *testing.T) {
	t.Run("AutoLinkNodes", testAutoLinkNodes)
	t.Run("ParagraphNodes", testParagraphNodes)
}

func testAutoLinkNodes(t *testing.T) {
	lexical.ResetNodes()
	lexical.RegisterNodes(&AutoLinkNode{}, &ParagraphNode{}, &TextNode{})
	message := `{"root":{"children":[{"children":[{"detail":0,"format":0,"mode":"normal","style":"","text":"with a link? ","type":"text","version":1},{"children":[{"detail":0,"format":0,"mode":"normal","style":"","text":"www.google.com","type":"text","version":1}],"direction":"ltr","format":"","indent":0,"type":"autolink","version":1,"rel":null,"target":null,"title":null,"url":"https://www.google.com","isUnlinked":true},{"detail":0,"format":0,"mode":"normal","style":"","text":" cool!","type":"text","version":1}],"direction":"ltr","format":"","indent":0,"type":"paragraph","version":1,"textFormat":0,"textStyle":""}],"direction":"ltr","format":"","indent":0,"type":"root","version":1}}`

	var root RootNode
	err := json.Unmarshal([]byte(message), &root)
	if err != nil {
		t.Fatal("json.Unmarshal err:", err)
	}

	if len(root.Root.Children) != 1 {
		t.Fatalf("expected length of root children 1; got %d", len(root.Root.Children))
	}

	paragraph, ok := root.Root.Children[0].(*ParagraphNode)
	if !ok {
		t.Fatal("root child is not a paragraph")
	}

	if len(paragraph.Children) != 3 {
		t.Fatalf("expected length of paragraph children 3; got %d", len(paragraph.Children))
	}

	autolink, ok := paragraph.Children[1].(*AutoLinkNode)
	if !ok {
		t.Fatal("paragraph child is not an autolink")
	}

	expectedURL := "https://www.google.com"
	if autolink.URL != expectedURL {
		t.Fatalf("expected url %s; got %s", expectedURL, autolink.URL)
	}

	if !autolink.IsUnlinked {
		t.Fatalf("expected isUnlinked %t; got %t", true, autolink.IsUnlinked)
	}
}

func testParagraphNodes(t *testing.T) {
	lexical.ResetNodes()
	lexical.RegisterNodes(&AutoLinkNode{}, &ParagraphNode{}, &TextNode{})
	message := `{"root":{"children":[{"children":[{"detail":0,"format":0,"mode":"normal","style":"","text":"with a link? ","type":"text","version":1},{"children":[{"detail":0,"format":0,"mode":"normal","style":"","text":"www.google.com","type":"text","version":1}],"direction":"ltr","format":"","indent":0,"type":"autolink","version":1,"rel":null,"target":null,"title":null,"url":"https://www.google.com","isUnlinked":true},{"detail":0,"format":0,"mode":"normal","style":"","text":" cool!","type":"text","version":1}],"direction":"ltr","format":"","indent":0,"type":"paragraph","version":1,"textFormat":1,"textStyle":"style"}],"direction":"ltr","format":"","indent":0,"type":"root","version":1}}`

	var root RootNode
	err := json.Unmarshal([]byte(message), &root)
	if err != nil {
		t.Fatal("json.Unmarshal err:", err)
	}

	if len(root.Root.Children) != 1 {
		t.Fatalf("expected length of root children 1; got %d", len(root.Root.Children))
	}

	paragraph, ok := root.Root.Children[0].(*ParagraphNode)
	if !ok {
		t.Fatal("root child is not a paragraph")
	}

	if len(paragraph.Children) != 3 {
		t.Fatalf("expected length of paragraph children 3; got %d", len(paragraph.Children))
	}

	expectedTextFormat := 1
	if paragraph.TextFormat != expectedTextFormat {
		t.Fatalf("expected textFormat %d; got %d", expectedTextFormat, paragraph.TextFormat)
	}

	expectedTextStyle := "style"
	if paragraph.TextStyle != expectedTextStyle {
		t.Fatalf("expected textStyle %s; got %s", expectedTextStyle, paragraph.TextStyle)
	}
}
