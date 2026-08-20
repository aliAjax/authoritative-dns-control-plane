package dns_wire

import "testing"

func TestEncodeParsePreservesQuestion(t *testing.T) {
	in := Message{Header: Header{ID: 7, Flags: 0x0100}, Questions: []Question{{Name: "www.example.com", Type: 1, Class: 1}}}
	out, err := Parse(Encode(in))
	if err != nil {
		t.Fatal(err)
	}
	if out.Header.ID != in.Header.ID || len(out.Questions) != 1 || out.QuestionName() != "www.example.com" {
		t.Fatalf("unexpected parsed message: %#v", out)
	}
}
