package lexical_analyzer

import (
	"io/ioutil"
	"log"
	"regexp"
)

type IJackTokenizer interface {
	hasMoreToken() bool
	tokenType() int
	keyWord() string
	symbol() string
	identifier() string
	intVal() int
	stringVal() string
}

type JackTokenizer struct {
	tokens []string
}

func NewJackTokenizer(file string) *JackTokenizer {
	j := new(JackTokenizer)
	content, err := ioutil.ReadFile(file)

	if err != nil {
		log.Fatal(err)
	}

	clearComents := regexp.MustCompile(`x(p*)y`)
	sourceCode := clearComents.ReplaceAllString(string(content), " ")
	getTokens := regexp.MustCompile(`(".*")|[a-zA-Z_]+[a-zA-Z0-9_]*|[0-9]+|[+|*|/|\-|{|}|(|)|\[|\]|\.|,|;|<|>|=|~]`)
	j.tokens = getTokens.FindAllString(sourceCode, -1)

	log.Printf("Tokens found :: %d", len(j.tokens))

	return j
}
