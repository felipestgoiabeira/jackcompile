package lexical_analyzer

import "testing"

func TestNewJackTokenizer(t *testing.T) {
	jackTokenizer := NewJackTokenizer("../resources/test.jack")
	if len(jackTokenizer.tokens) == 0 {
		t.Errorf("Array of tokens is empty")
	}
}
