package parser

import (
	"errors"
	"fmt"
	la "jackcompile/lexical_analysis"
	"jackcompile/utils"
	"regexp"
)

func eat(expectedToken string, jackTokenizer *la.JackTokenizer) string {
	jackTokenizer.Advance()
	token := jackTokenizer.GetCurToken()
	if expectedToken != token.GetToken() {
		panic(errors.New("Not expected token :: " + token.GetToken()))
	}
	return XmlToken(token)
}

func isExpectedToken(expectedToken string, token la.Token) {
	if expectedToken != token.GetToken() {
		panic(errors.New(expectedToken + " expected, got " + token.GetToken()))
	}
}

func isOptionalExpectedToken(expectedToken string, token la.Token) bool {
	if expectedToken == token.GetToken() {
		return true
	}
	return false
}

func isOperator(op la.Token) bool {
	token := op.GetToken()
	isOperator, _ := regexp.MatchString(
		`\+|-|\*|\/|&|\||<|>|=`,
		token,
	)
	return isOperator
}

func XmlToken(token la.Token) string {
	xmlSpecialCharacters := map[string]string{
		"<":  "&lt;",
		">":  "&gt;",
		"&":  "&amp;",
		"'":  "&#39",
		"\"": "&#34",
	}
	_, isSpecialCharacter := xmlSpecialCharacters[token.GetToken()]
	if isSpecialCharacter {
		token.SetToken(xmlSpecialCharacters[token.GetToken()])
	}
	return fmt.Sprintf("<%s> %s </%s>", token.GetType(), token.GetToken(), token.GetType())
}

func GetTokenTypeTerms() []la.TokenType {
	return []la.TokenType{la.IDENTIFIER, la.INTEGER, la.STRING}
}

func CompileTerm(token la.Token) []string {
	if !utils.ContainsTokenType(GetTokenTypeTerms(), token.GetType()) {
		panic(errors.New("A term expected was not found, token found :: " + token.GetToken()))
	}
	var result []string
	result = append(result, "<term>")
	result = append(result, "  "+XmlToken(token))
	result = append(result, "</term>")
	return result
}

func CompileExpression(jackTokenizer *la.JackTokenizer) []string {
	var result []string
	result = append(result, "<expression>")

	for _, term := range CompileTerm(jackTokenizer.GetCurToken()) {
		result = append(result, "  "+term)
	}

	for jackTokenizer.HasPeekToken() && isOperator(jackTokenizer.GetPeekToken()) {
		jackTokenizer.Advance()
		result = append(result, "  "+XmlToken(jackTokenizer.GetCurToken()))
		jackTokenizer.Advance()
		for _, term := range CompileTerm(jackTokenizer.GetCurToken()) {
			result = append(result, "  "+term)
		}
	}

	result = append(result, "</expression>")
	return result
}

func CompileIfStatement(jackTokenizer *la.JackTokenizer) []string {
	var result []string
	isExpectedToken("if", jackTokenizer.GetCurToken())
	result = append(result, "<ifStatement>")
	result = utils.AppendIndent(result, XmlToken(jackTokenizer.GetCurToken()))
	result = utils.AppendIndent(result, eat("(", jackTokenizer))
	jackTokenizer.Advance()
	result = utils.AppendIndent(result, CompileExpression(jackTokenizer)...)
	result = utils.AppendIndent(result, eat(")", jackTokenizer))
	result = utils.AppendIndent(result, eat("{", jackTokenizer))
	result = utils.AppendIndent(result, eat("}", jackTokenizer))
	result = append(result, "</ifStatement>")
	return result
}

func CompileLetStatement(jackTokenizer *la.JackTokenizer) []string {
	var result []string
	isExpectedToken("let", jackTokenizer.GetCurToken())
	result = append(result, "<letStatement>")
	result = utils.AppendIndent(result, XmlToken(jackTokenizer.GetCurToken()))
	jackTokenizer.Advance()
	result = utils.AppendIndent(result, XmlToken(jackTokenizer.GetCurToken()))
	result = utils.AppendIndent(result, eat("=", jackTokenizer))
	jackTokenizer.Advance()
	result = utils.AppendIndent(result, CompileExpression(jackTokenizer)...)
	result = utils.AppendIndent(result, eat(";", jackTokenizer))
	result = append(result, "</letStatement>")
	return result
}

func CompileReturnStatement(jackTokenizer *la.JackTokenizer) []string {
	var result []string
	isExpectedToken("return", jackTokenizer.GetCurToken())
	result = append(result, "<returnStatement>")
	result = utils.AppendIndent(result, XmlToken(jackTokenizer.GetCurToken()))

	peekToken := jackTokenizer.GetPeekToken()
	if peekToken.GetToken() == ";" {
		result = utils.AppendIndent(result, eat(";", jackTokenizer))
		result = append(result, "</returnStatement>")
		return result
	}

	jackTokenizer.Advance()

	result = utils.AppendIndent(result, CompileExpression(jackTokenizer)...)
	result = utils.AppendIndent(result, eat(";", jackTokenizer))
	result = append(result, "</returnStatement>")

	return result
}

func CompileExpressionList(jackTokenizer *la.JackTokenizer) []string {
	var result []string
	result = append(result, "<expressionList>")
	result = utils.AppendIndent(result, CompileExpression(jackTokenizer)...)
	for jackTokenizer.HasPeekToken() && isOptionalExpectedToken(",", jackTokenizer.GetPeekToken()) {
		result = utils.AppendIndent(result, eat(",", jackTokenizer))
		jackTokenizer.Advance()
		result = utils.AppendIndent(result, CompileExpression(jackTokenizer)...)
	}
	result = append(result, "</expressionList>")
	return result
}
