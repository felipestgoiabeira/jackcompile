package lexical_analysis

import (
	"log"
	"regexp"
	"strings"
)

type Token struct {
	tokenType TokenType
	token     string
}

func NewToken(token string) *Token {
	t := new(Token)
	t.token = token

	t.isKeyword()
	t.isSymbol()
	t.isIntegerConst()
	t.isStringConst()
	t.isIdentifier()

	return t
}

func (t *Token) isKeyword() {
	if t.tokenType != "" {
		return
	}

	isKeyword, _ := regexp.MatchString(
		`int|class|constructor|function|method|field|static|var|char|boolean|void|true|false|null|this|let|do|if|else|while|return`,
		t.token,
	)

	if isKeyword {
		log.Printf("The token %s is a keyword", t.token)
		t.tokenType = KEYWORD
	}
}

func (t *Token) isSymbol() {
	if t.tokenType != "" {
		return
	}

	if strings.Contains("{}()[].,;+-*/&|<>=~", t.token) {
		log.Printf("The token %s is a symbol", t.token)
		t.tokenType = SYMBOL
	}
}

func (t *Token) isIntegerConst() {
	if t.tokenType != "" {
		return
	}

	integer, _ := regexp.MatchString(`[0-9]+`, t.token)
	if integer {
		log.Printf("The token %s is a integerConst", t.token)
		t.tokenType = INTEGER
	}
}

func (t *Token) isStringConst() {
	if t.tokenType != "" {
		return
	}

	isString, _ := regexp.MatchString(`".*"`, t.token)
	if isString {
		log.Printf("The token %s is a stringConst", t.token)
		r := regexp.MustCompile(`"`)
		t.token = r.ReplaceAllString(t.token, "")
		t.tokenType = STRING
	}
}

func (t *Token) isIdentifier() {
	if t.tokenType == "" {
		log.Printf("The token %s is a identifier", t.token)
		t.tokenType = IDENTIFIER
	}
}
