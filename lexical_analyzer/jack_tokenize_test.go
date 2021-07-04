package lexical_analyzer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

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

func TestNewJackTokenizerMustBuildTheExpectedXMLTree(t *testing.T) {
	jackTokenizer := NewJackTokenizer("../resources/complex_code.jack")
	content, err := ioutil.ReadFile("../resources/xml/complex_code.xml")

	if err != nil {
		log.Fatal(err)
	}

	expectedResult := string(content)

	fmt.Println(expectedResult)

	xmlSpecialCharacters := map[string]string{
		"<":  "&lt",
		">":  "&gt",
		"&":  "&#38",
		"'":  "&#39",
		"\"": "&#34",
	}

	var resultBuffer bytes.Buffer

	resultBuffer.WriteString("<tokens>")

	for _, token := range jackTokenizer.tokens {
		_, isSpecialCharacter := xmlSpecialCharacters[token.token]
		if isSpecialCharacter {
			token.token = xmlSpecialCharacters[token.token]
		}
		resultBuffer.WriteString(fmt.Sprintf("\n<%s> %s </%s>", token.tokenType, token.token, token.tokenType))
	}

	resultBuffer.WriteString("\n</tokens>")

	result := resultBuffer.String()

	if expectedResult != result {
		t.Errorf("The expected result don't match with the result")
	}
}
