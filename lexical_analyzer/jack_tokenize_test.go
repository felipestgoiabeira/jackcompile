package lexical_analyzer

import "testing"

func TestNewJackTokenizerMustReturnANonEmptySlice(t *testing.T) {
	jackTokenizer := NewJackTokenizer("../resources/simple_code.jack")
	if len(jackTokenizer.tokens) == 0 {
		t.Errorf("Array of tokens is empty")
	}
}

func TestNewJackTokenizerMustReturnAEmptySlice(t *testing.T) {
	jackTokenizer := NewJackTokenizer("../resources/only_comments.jack")
	if len(jackTokenizer.tokens) > 0 {
		t.Errorf("Array of tokens is not empty")
	}
}
