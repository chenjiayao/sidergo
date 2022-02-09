package border

import (
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
