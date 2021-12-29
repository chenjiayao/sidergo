package helper

import (
	"testing"
)

func TestBbyteToSString(t *testing.T) {

	b := [][]byte{
		[]byte("go"),
		[]byte("redis"),
		[]byte("training"),
	}

	got := BbyteToSString(b)

	if len(got) != len(b) {
		t.Errorf("len(BbyteToSString(b)) = %d, want = %d", len(got), len(b))
		return
	}

	for i := 0; i < len(got); i++ {
		if string(b[i]) != got[i] {
			t.Errorf("BbyteToSString(b)[%d] = %s, want = %s", i, got[i], string(b[i]))
		}
	}

}
