package lexical_analysis

import "testing"

func TestIntegerToken(t *testing.T) {
	token := NewToken("15")
	if token.tokenType != INTEGER {
		t.Fatal("Token isn't a integer")
	}
}

func TestStringToken(t *testing.T) {
	token := NewToken("\"string\"")

	if token.tokenType != STRING {
		t.Fatal("Token isn't a integer")
	}

	if token.token != "string" {
		t.Fatal("Must return only the string")
	}
}

func TestSymbolToken(t *testing.T) {
	token := NewToken("~")

	if token.tokenType != SYMBOL {
		t.Fatal("Token isn't a symbol")
	}
}

func TestKeywordToken(t *testing.T) {
	token := NewToken("return")

	if token.tokenType != KEYWORD {
		t.Fatal("Token isn't a keyword")
	}
}

func TestIdentifierToken(t *testing.T) {
	token := NewToken("aleatory")

	if token.tokenType != IDENTIFIER {
		t.Fatal("Token isn't a identifier")
	}
}
