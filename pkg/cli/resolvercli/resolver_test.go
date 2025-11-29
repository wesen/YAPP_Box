package resolvercli

import "testing"

func TestDecodeDocumentWithComments(t *testing.T) {
	yaml := `
pcb:
  length: 50.0 # board length
vars:
  spacing: 10 # spacing comment
`
	doc, comments, err := decodeDocumentWithComments([]byte(yaml))
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if doc == nil {
		t.Fatalf("expected doc")
	}
	got := comments["pcb.length"]
	if len(got) == 0 || got[0] != "board length" {
		t.Fatalf("expected board length comment, got %#v", got)
	}
	got = comments["vars.spacing"]
	if len(got) == 0 || got[0] != "spacing comment" {
		t.Fatalf("expected spacing comment, got %#v", got)
	}
}
