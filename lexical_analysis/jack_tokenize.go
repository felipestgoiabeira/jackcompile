package lexical_analysis

import (
	"io/ioutil"
	"log"
	"regexp"
)

type JackTokenizer struct {
	tokens []Token
	pos    int32
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

func (j *JackTokenizer) Advance() {
	j.pos++
}

func (j *JackTokenizer) HasMoreTokens() bool {
	posCur := j.pos + 1
	return !(int(posCur) > len(j.tokens))
}

func (j *JackTokenizer) HasPeekToken() bool {
	posCur := j.pos + 2
	return !(int(posCur) > len(j.tokens))
}

func (j *JackTokenizer) GetCurToken() Token {
	return j.tokens[j.pos]
}

func (j *JackTokenizer) GetCurTokenType() TokenType {
	return j.tokens[j.pos].GetType()
}

func (j *JackTokenizer) GetPeekToken() Token {
	return j.tokens[j.pos+1]
}

func (j *JackTokenizer) GetPeekTokenType() TokenType {
	return j.tokens[j.pos+1].GetType()
}
