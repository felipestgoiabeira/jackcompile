package lexical_analyzer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	jackTokenizer := NewJackTokenizer("../resources/Square.jack")
	content, err := ioutil.ReadFile("../resources/xml/SquareT.xml")

	if err != nil {
		log.Fatal(err)
	}

	expectedResult := string(content)

	xmlSpecialCharacters := map[string]string{
		"<":  "&lt;",
		">":  "&gt;",
		"&":  "&amp;",
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

	writeResultToFile(resultBuffer)

	if expectedResult != result {
		t.Errorf("The expected result don't match with the result")
	}
}

func writeResultToFile(resultBuffer bytes.Buffer) {
	f, err := os.Create("../resources/tests/results/SquareResult.xml")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.Write(resultBuffer.Bytes())

	if err2 != nil {
		log.Fatal(err2)
	}
}
