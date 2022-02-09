package border

import (
	"math"
	"testing"
)

func TestParserBorder(t *testing.T) {

	s1 := "12.3"
	border1, _ := ParserBorder(s1)
	if !(border1.Include && border1.Inf == 0 && border1.Value == 12.3) {
		t.Errorf("parase %s failed", s1)
	}

	s2 := "(12.3"
	border2, _ := ParserBorder(s2)
	if !(!border2.Include && border2.Inf == 0 && border2.Value == 12.3) {
		t.Errorf("parase %s failed", s2)
	}

	s3 := "-inf"
	border3, _ := ParserBorder(s3)
	if !(border3.Include && border3.Inf == -1 && border3.Value == 0) {
		t.Errorf("parase %s failed", s3)
	}

	s4 := "(-inf"
	border4, _ := ParserBorder(s4)
	if !(border4.Include && border4.Inf == -1) {
		t.Errorf("parase %s failed", s4)
	}
}

func TestBorder_Greater(t *testing.T) {
	b1 := &Border{
		Inf: negativeInf,
	}

	if b1.Greater(math.Inf(-1)) {
		t.Errorf("b1.Greater(math.Inf(-1)) failed")
	}

	if b1.Greater(math.Inf(1)) {
		t.Errorf("b1.Greater(math.Inf(1)) failed")
	}

	////
	b2 := &Border{
		Inf:     0,
		Include: true,
		Value:   12,
	}

	if !b2.Greater(12) {
		t.Errorf("b2.Greater(math.Inf(12)) failed")
	}
	if b2.Greater(13) {
		t.Errorf("b2.Greater(math.Inf(13)) failed")
	}

	if !b2.Greater(11) {
		t.Errorf("b2.Greater(math.Inf(11)) failed")
	}

	///
	b3 := &Border{
		Inf:     0,
		Include: false,
		Value:   12,
	}

	if b3.Greater(12) {
		t.Errorf("b2.Greater(math.Inf(12)) failed")
	}
	if b3.Greater(13) {
		t.Errorf("b3.Greater(math.Inf(13)) failed")
	}

	if !b3.Greater(11) {
		t.Errorf("b3.Greater(math.Inf(11)) failed")
	}
}
