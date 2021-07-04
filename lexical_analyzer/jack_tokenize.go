package lexical_analyzer

import (
	"io/ioutil"
	"log"
	"regexp"
)

type JackTokenizer struct {
	tokens []Token
}

func NewJackTokenizer(file string) *JackTokenizer {
	j := new(JackTokenizer)
	content, err := ioutil.ReadFile(file)

	if err != nil {
		log.Fatal(err)
	}

	clearComents := regexp.MustCompile(`(//.*\n)|(/\*(.|\n)*?\*/)`)
	sourceCode := clearComents.ReplaceAllString(string(content), " ")
	getTokens := regexp.MustCompile(`(".*")|[a-zA-Z_]+[a-zA-Z0-9_]*|[0-9]+|[+|*|/|\-|{|}|(|)|\[|\]|\.|,|;|<|>|=|~|&]`)

	for _, sToken := range getTokens.FindAllString(sourceCode, -1) {
		j.tokens = append(j.tokens, *NewToken(sToken))
	}

	log.Printf("Tokens found :: %d", len(j.tokens))

	return j
}
