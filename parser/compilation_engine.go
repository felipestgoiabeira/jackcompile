package parser

import (
	"errors"
	"fmt"
	la "jackcompile/lexical_analysis"
	"jackcompile/utils"
	"regexp"
)

func CompileStatement(statement la.Token, jackTokenize *la.JackTokenizer) []string {
	var result []string
	statements := map[string]func(*la.JackTokenizer) []string{
		"let":    CompileLetStatement,
		"do":     CompileDoStatement,
		"if":     CompileIfStatement,
		"return": CompileReturnStatement,
		"while":  CompileWhileStatement,
	}
	compileStatement, exists := statements[statement.GetToken()]
	if exists {
		result = utils.AppendIndent(result, compileStatement(jackTokenize)...)
		return result

	}
	panic(errors.New("Not a statement"))

}

func eat(expectedToken string, jackTokenizer *la.JackTokenizer) string {
	jackTokenizer.Advance()
	token := jackTokenizer.GetCurToken()
	if expectedToken != token.GetToken() {
		panic(errors.New(expectedToken + " expected, got " + token.GetToken()))
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

func isOneOf(regex string, op la.Token) bool {
	token := op.GetToken()
	isOperator, _ := regexp.MatchString(
		regex,
		token,
	)
	return isOperator
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
	result = utils.AppendIndent(result, compileStatements(jackTokenizer)...)
	result = utils.AppendIndent(result, eat("}", jackTokenizer))
	if isOptionalExpectedToken("else", jackTokenizer.GetPeekToken()) {
		jackTokenizer.Advance()
		result = utils.AppendIndent(result, XmlToken(jackTokenizer.GetCurToken()))
		result = utils.AppendIndent(result, eat("{", jackTokenizer))
		result = utils.AppendIndent(result, compileStatements(jackTokenizer)...)
		result = utils.AppendIndent(result, eat("}", jackTokenizer))
	}
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
	curToken := jackTokenizer.GetCurToken()
	if curToken.GetToken() == ")" {
		result = append(result, "</expressionList>")
		return result
	}
	result = utils.AppendIndent(result, CompileExpression(jackTokenizer)...)
	for jackTokenizer.HasPeekToken() && isOptionalExpectedToken(",", jackTokenizer.GetPeekToken()) {
		result = utils.AppendIndent(result, eat(",", jackTokenizer))
		jackTokenizer.Advance()
		result = utils.AppendIndent(result, CompileExpression(jackTokenizer)...)
	}
	result = append(result, "</expressionList>")
	return result
}

func compileExpressionListInsideParentheses(jackTokenizer *la.JackTokenizer) []string {
	var result []string
	result = utils.AppendIndent(result, eat("(", jackTokenizer))
	jackTokenizer.Advance()
	result = utils.AppendIndent(result, CompileExpressionList(jackTokenizer)...)
	curToken := jackTokenizer.GetCurToken()
	if curToken.GetToken() == ")" {
		result = utils.AppendIndent(result, XmlToken(curToken))
		return result
	}
	result = utils.AppendIndent(result, eat(")", jackTokenizer))
	return result
}

func CompileDoStatement(jackTokenizer *la.JackTokenizer) []string {
	var result []string
	isExpectedToken("do", jackTokenizer.GetCurToken())
	result = append(result, "<doStatement>")
	result = utils.AppendIndent(result, XmlToken(jackTokenizer.GetCurToken()))
	jackTokenizer.Advance()
	result = utils.AppendIndent(result, XmlToken(jackTokenizer.GetCurToken()))
	peekToken := jackTokenizer.GetPeekToken()
	if peekToken.GetToken() == "(" {
		result = append(result, compileExpressionListInsideParentheses(jackTokenizer)...)
		result = utils.AppendIndent(result, eat(";", jackTokenizer))
		result = append(result, "</doStatement>")
		return result
	}
	result = utils.AppendIndent(result, eat(".", jackTokenizer))
	jackTokenizer.Advance()
	result = utils.AppendIndent(result, XmlToken(jackTokenizer.GetCurToken()))
	result = append(result, compileExpressionListInsideParentheses(jackTokenizer)...)
	result = utils.AppendIndent(result, eat(";", jackTokenizer))
	result = append(result, "</doStatement>")
	return result
}

func CompileWhileStatement(jackTokenizer *la.JackTokenizer) []string {
	var result []string
	isExpectedToken("while", jackTokenizer.GetCurToken())
	result = append(result, "<whileStatement>")
	result = utils.AppendIndent(result, XmlToken(jackTokenizer.GetCurToken()))
	result = utils.AppendIndent(result, eat("(", jackTokenizer))
	jackTokenizer.Advance()
	result = utils.AppendIndent(result, CompileExpression(jackTokenizer)...)

	result = utils.AppendIndent(result, eat(")", jackTokenizer))
	result = utils.AppendIndent(result, eat("{", jackTokenizer))
	for jackTokenizer.HasPeekToken() && !isOptionalExpectedToken("}", jackTokenizer.GetPeekToken()) {
		jackTokenizer.Advance()
		result = utils.AppendIndent(result, CompileStatement(jackTokenizer.GetCurToken(), jackTokenizer)...)
	}
	result = utils.AppendIndent(result, eat("}", jackTokenizer))
	result = append(result, "</whileStatement>")
	return result
}

func CompileVarDeclaration(jackTokenizer *la.JackTokenizer) []string {
	var result []string
	isExpectedToken("var", jackTokenizer.GetCurToken())
	result = append(result, "<varDec>")
	result = utils.AppendIndent(result, XmlToken(jackTokenizer.GetCurToken()))
	jackTokenizer.Advance()
	result = utils.AppendIndent(result, XmlToken(jackTokenizer.GetCurToken()))
	jackTokenizer.Advance()
	result = utils.AppendIndent(result, XmlToken(jackTokenizer.GetCurToken()))
	result = utils.AppendIndent(result, eat(";", jackTokenizer))
	result = append(result, "</varDec>")
	return result
}

func CompileParameterList(jackTokenizer *la.JackTokenizer) []string {
	var result []string
	result = append(result, "<parameterList>")
	peekToken := jackTokenizer.GetPeekToken()
	if peekToken.GetToken() == ")" {
		result = append(result, "</parameterList>")
		return result
	}
	result = utils.AppendIndent(result, XmlToken(jackTokenizer.GetCurToken()))
	for jackTokenizer.HasPeekToken() && isOptionalExpectedToken(",", jackTokenizer.GetPeekToken()) {
		jackTokenizer.Advance()
		result = utils.AppendIndent(result, XmlToken(jackTokenizer.GetCurToken()))
		jackTokenizer.Advance()
		result = utils.AppendIndent(result, XmlToken(jackTokenizer.GetCurToken()))
	}
	result = append(result, "</parameterList>")
	return result
}

func compileStatements(jackTokenizer *la.JackTokenizer) []string {
	var result []string
	result = append(result, "<statements>")
	for jackTokenizer.HasPeekToken() && !isOptionalExpectedToken("}", jackTokenizer.GetPeekToken()) {
		jackTokenizer.Advance()
		result = append(result, CompileStatement(jackTokenizer.GetCurToken(), jackTokenizer)...)
	}
	result = append(result, "</statements>")
	return result
}

func compileSubroutineBody(jackTokenizer *la.JackTokenizer) []string {
	var result []string
	result = append(result, "<subroutineBody>")
	result = utils.AppendIndent(result, eat("{", jackTokenizer))
	for jackTokenizer.HasPeekToken() && isOptionalExpectedToken("var", jackTokenizer.GetPeekToken()) {
		jackTokenizer.Advance()
		result = utils.AppendIndent(result, CompileVarDeclaration(jackTokenizer)...)
	}
	result = utils.AppendIndent(result, compileStatements(jackTokenizer)...)
	result = utils.AppendIndent(result, eat("}", jackTokenizer))
	result = append(result, "</subroutineBody>")
	return result
}

func CompileSubroutine(jackTokenizer *la.JackTokenizer) []string {
	if !isOneOf(`constructor|function|method`, jackTokenizer.GetCurToken()) {
		panic(errors.New("Not a subroutine!"))
	}
	var result []string
	result = append(result, "<subroutineDec>")
	result = utils.AppendIndent(result, XmlToken(jackTokenizer.GetCurToken()))
	jackTokenizer.Advance()
	result = utils.AppendIndent(result, XmlToken(jackTokenizer.GetCurToken()))
	jackTokenizer.Advance()
	if isOptionalExpectedToken("new", jackTokenizer.GetCurToken()) {
		result = utils.AppendIndent(result, XmlToken(jackTokenizer.GetCurToken()))
		jackTokenizer.Advance()
	}
	result = utils.AppendIndent(result, XmlToken(jackTokenizer.GetCurToken()))
	result = utils.AppendIndent(result, eat("(", jackTokenizer))
	result = utils.AppendIndent(result, CompileParameterList(jackTokenizer)...)
	result = utils.AppendIndent(result, eat(")", jackTokenizer))
	result = utils.AppendIndent(result, compileSubroutineBody(jackTokenizer)...)
	result = append(result, "</subroutineDec>")
	return result
}

func CompileClassVarDec(jackTokenizer *la.JackTokenizer) []string {
	if !isOneOf(`static|field`, jackTokenizer.GetCurToken()) {
		panic(errors.New("Not a subroutine!"))
	}
	var result []string
	result = append(result, "<classVarDec>")
	result = utils.AppendIndent(result, XmlToken(jackTokenizer.GetCurToken()))
	jackTokenizer.Advance()
	result = utils.AppendIndent(result, XmlToken(jackTokenizer.GetCurToken()))
	jackTokenizer.Advance()
	result = utils.AppendIndent(result, XmlToken(jackTokenizer.GetCurToken()))
	result = utils.AppendIndent(result, eat(";", jackTokenizer))
	result = append(result, "</classVarDec>")
	return result
}

func compileDeclaration(declaration la.Token, jackTokenize *la.JackTokenizer) ([]string, bool) {

	var result []string
	declarations := map[string]func(*la.JackTokenizer) []string{
		"field":       CompileClassVarDec,
		"static":      CompileClassVarDec,
		"function":    CompileSubroutine,
		"constructor": CompileSubroutine,
		"method":      CompileSubroutine,
	}

	compileDeclaration, valid := declarations[declaration.GetToken()]
	if valid {
		result = utils.AppendIndent(result, compileDeclaration(jackTokenize)...)
	}
	return result, valid
}

func CompileClass(jackTokenizer *la.JackTokenizer) []string {
	isExpectedToken("class", jackTokenizer.GetCurToken())
	var result []string
	result = append(result, "<class>")
	result = utils.AppendIndent(result, XmlToken(jackTokenizer.GetCurToken()))
	jackTokenizer.Advance()
	result = utils.AppendIndent(result, XmlToken(jackTokenizer.GetCurToken()))
	result = utils.AppendIndent(result, eat("{", jackTokenizer))
	totalCompiled := false
	for !totalCompiled {
		jackTokenizer.Advance()
		compiled, valid := compileDeclaration(jackTokenizer.GetCurToken(), jackTokenizer)
		if !valid {
			totalCompiled = true
		}
		result = append(result, compiled...)
	}
	result = utils.AppendIndent(result, XmlToken(jackTokenizer.GetCurToken()))
	result = append(result, "</class>")
	return result
}
