package messages

import "testing"

func TestPlainTextFromMarkdown(t *testing.T) {
	input := "Hello **Bob**\n\n```go\nfmt.Println(\"x\")\n```"
	got := PlainTextFromMarkdown(input)
	want := "Hello Bob fmt.Println(\"x\")"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
