package cli

import (
	"io"
	"testing"
)

// chunks returns its byte slices one Read at a time, as a terminal would.
type chunks [][]byte

func (c *chunks) Read(p []byte) (int, error) {
	if len(*c) == 0 {
		return 0, io.EOF
	}
	n := copy(p, (*c)[0])
	*c = (*c)[1:]
	return n, nil
}

type got struct {
	key  int
	char rune
}

func keys(t *testing.T, reads ...string) []got {
	t.Helper()
	pending = nil
	in := chunks{}
	for _, r := range reads {
		in = append(in, []byte(r))
	}
	var out []got
	for {
		key, char, err := readKey(&in)
		if err == io.EOF {
			return out
		}
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, got{key, char})
	}
}

func equal(a, b []got) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestReadKeyKeepsEveryKeyOfAReadOrPaste(t *testing.T) {
	want := []got{{keyUnknown, 'a'}, {keyUnknown, 'b'}, {keySpace, ' '}, {keyUnknown, 'c'}, {keyEnter, '\n'}}
	if out := keys(t, "ab c\r"); !equal(out, want) {
		t.Fatalf("one read: %v", out)
	}
	if out := keys(t, "a", "b", " ", "c", "\r"); !equal(out, want) {
		t.Fatalf("separate reads: %v", out)
	}
}

func TestReadKeyParsesSequencesBetweenCharacters(t *testing.T) {
	out := keys(t, "x\x1b[Ay\x1bOB\x1b[3~\x1b[1;5C\x1b[13;2u\x7f\b")
	want := []got{
		{keyUnknown, 'x'}, {keyUp, 0}, {keyUnknown, 'y'}, {keyDown, 0}, {keyDelete, 0},
		{keyRight, 0}, {keyShiftEnter, '\n'}, {keyBackspace, 0}, {keyBackspace, 0},
	}
	if !equal(out, want) {
		t.Fatalf("got %v", out)
	}
}

func TestReadKeyJoinsSplitCharactersAndSequences(t *testing.T) {
	// "ü" is two bytes and "→" three; both, and an arrow key, arrive split.
	out := keys(t, "\xc3", "\xbc\xe2\x86", "\x92", "\x1b[", "D")
	want := []got{{keyUnknown, 'ü'}, {keyUnknown, '→'}, {keyLeft, 0}}
	if !equal(out, want) {
		t.Fatalf("got %v", out)
	}
}

func TestReadKeyTreatsALoneEscapeAsTheEscapeKey(t *testing.T) {
	out := keys(t, "\x1b", "q")
	want := []got{{keyEscape, 0}, {keyUnknown, 'q'}}
	if !equal(out, want) {
		t.Fatalf("got %v", out)
	}
}
