package parser

import (
	"errors"
	"fmt"
	la "jackcompile/lexical_analysis"
	"jackcompile/utils"
	"regexp"
)

func CompileStatement(statement la.Token, jackTokenize *la.JackTokenizer) []string {
	var xmlResult []string
	statements := map[string]func(*la.JackTokenizer) []string{
		"let":    CompileLetStatement,
		"do":     CompileDoStatement,
		"if":     CompileIfStatement,
		"return": CompileReturnStatement,
		"while":  CompileWhileStatement,
	}
	compileStatement, exists := statements[statement.GetToken()]
	if exists {
		utils.AppendIndent(&xmlResult, compileStatement(jackTokenize)...)
		return xmlResult

	}
	panic(errors.New("Not a statement"))

}

func eat(expectedToken string, jackTokenizer *la.JackTokenizer, xmlResult *[]string) {
	jackTokenizer.Advance()
	token := jackTokenizer.GetCurToken()
	if expectedToken != token.GetToken() {
		panic(errors.New(expectedToken + " expected, got " + token.GetToken()))
	}
	utils.AppendIndent(xmlResult, xmlToken(token))
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

func xmlToken(token la.Token) string {
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
	var xmlResult []string
	xmlResult = append(xmlResult, "<term>")
	xmlResult = append(xmlResult, "  "+xmlToken(token))
	xmlResult = append(xmlResult, "</term>")
	return xmlResult
}

func CompileExpression(jackTokenizer *la.JackTokenizer) []string {
	var xmlResult []string
	xmlResult = append(xmlResult, "<expression>")

	for _, term := range CompileTerm(jackTokenizer.GetCurToken()) {
		xmlResult = append(xmlResult, "  "+term)
	}

	for jackTokenizer.HasPeekToken() && isOperator(jackTokenizer.GetPeekToken()) {
		jackTokenizer.Advance()
		xmlResult = append(xmlResult, "  "+xmlToken(jackTokenizer.GetCurToken()))
		jackTokenizer.Advance()
		for _, term := range CompileTerm(jackTokenizer.GetCurToken()) {
			xmlResult = append(xmlResult, "  "+term)
		}
	}

	xmlResult = append(xmlResult, "</expression>")
	return xmlResult
}

func xmlTokenAppendIndent(xmlResult *[]string, jackTokenizer *la.JackTokenizer) {
	utils.AppendIndent(xmlResult, xmlToken(jackTokenizer.GetCurToken()))
}

func CompileIfStatement(jackTokenizer *la.JackTokenizer) []string {
	var xmlResult []string
	isExpectedToken("if", jackTokenizer.GetCurToken())
	xmlResult = append(xmlResult, "<ifStatement>")
	utils.AppendIndent(&xmlResult, xmlToken(jackTokenizer.GetCurToken()))
	eat("(", jackTokenizer, &xmlResult)
	jackTokenizer.Advance()
	utils.AppendIndent(&xmlResult, CompileExpression(jackTokenizer)...)
	eat(")", jackTokenizer, &xmlResult)
	eat("{", jackTokenizer, &xmlResult)
	utils.AppendIndent(&xmlResult, compileStatements(jackTokenizer)...)
	eat("}", jackTokenizer, &xmlResult)
	if isOptionalExpectedToken("else", jackTokenizer.GetPeekToken()) {
		jackTokenizer.Advance()
		utils.AppendIndent(&xmlResult, xmlToken(jackTokenizer.GetCurToken()))
		eat("{", jackTokenizer, &xmlResult)
		utils.AppendIndent(&xmlResult, compileStatements(jackTokenizer)...)
		eat("}", jackTokenizer, &xmlResult)
	}
	xmlResult = append(xmlResult, "</ifStatement>")
	return xmlResult
}

func CompileLetStatement(jackTokenizer *la.JackTokenizer) []string {
	var xmlResult []string
	isExpectedToken("let", jackTokenizer.GetCurToken())
	xmlResult = append(xmlResult, "<letStatement>")
	utils.AppendIndent(&xmlResult, xmlToken(jackTokenizer.GetCurToken()))
	jackTokenizer.Advance()
	utils.AppendIndent(&xmlResult, xmlToken(jackTokenizer.GetCurToken()))
	eat("=", jackTokenizer, &xmlResult)
	jackTokenizer.Advance()
	utils.AppendIndent(&xmlResult, CompileExpression(jackTokenizer)...)
	eat(";", jackTokenizer, &xmlResult)
	xmlResult = append(xmlResult, "</letStatement>")
	return xmlResult
}

func CompileReturnStatement(jackTokenizer *la.JackTokenizer) []string {
	var xmlResult []string
	isExpectedToken("return", jackTokenizer.GetCurToken())
	xmlResult = append(xmlResult, "<returnStatement>")
	utils.AppendIndent(&xmlResult, xmlToken(jackTokenizer.GetCurToken()))

	peekToken := jackTokenizer.GetPeekToken()
	if peekToken.GetToken() == ";" {
		eat(";", jackTokenizer, &xmlResult)
		xmlResult = append(xmlResult, "</returnStatement>")
		return xmlResult
	}

	jackTokenizer.Advance()

	utils.AppendIndent(&xmlResult, CompileExpression(jackTokenizer)...)
	eat(";", jackTokenizer, &xmlResult)
	xmlResult = append(xmlResult, "</returnStatement>")

	return xmlResult
}

func CompileExpressionList(jackTokenizer *la.JackTokenizer) []string {
	var xmlResult []string
	xmlResult = append(xmlResult, "<expressionList>")
	curToken := jackTokenizer.GetCurToken()
	if curToken.GetToken() == ")" {
		xmlResult = append(xmlResult, "</expressionList>")
		return xmlResult
	}
	utils.AppendIndent(&xmlResult, CompileExpression(jackTokenizer)...)
	for jackTokenizer.HasPeekToken() && isOptionalExpectedToken(",", jackTokenizer.GetPeekToken()) {
		eat(",", jackTokenizer, &xmlResult)
		jackTokenizer.Advance()
		utils.AppendIndent(&xmlResult, CompileExpression(jackTokenizer)...)
	}
	xmlResult = append(xmlResult, "</expressionList>")
	return xmlResult
}

func compileExpressionListInsideParentheses(jackTokenizer *la.JackTokenizer) []string {
	var xmlResult []string
	eat("(", jackTokenizer, &xmlResult)
	jackTokenizer.Advance()
	utils.AppendIndent(&xmlResult, CompileExpressionList(jackTokenizer)...)
	curToken := jackTokenizer.GetCurToken()
	if curToken.GetToken() == ")" {
		utils.AppendIndent(&xmlResult, xmlToken(curToken))
		return xmlResult
	}
	eat(")", jackTokenizer, &xmlResult)
	return xmlResult
}

func CompileDoStatement(jackTokenizer *la.JackTokenizer) []string {
	var xmlResult []string
	isExpectedToken("do", jackTokenizer.GetCurToken())
	xmlResult = append(xmlResult, "<doStatement>")
	utils.AppendIndent(&xmlResult, xmlToken(jackTokenizer.GetCurToken()))
	jackTokenizer.Advance()
	utils.AppendIndent(&xmlResult, xmlToken(jackTokenizer.GetCurToken()))
	peekToken := jackTokenizer.GetPeekToken()
	if peekToken.GetToken() == "(" {
		xmlResult = append(xmlResult, compileExpressionListInsideParentheses(jackTokenizer)...)
		eat(";", jackTokenizer, &xmlResult)
		xmlResult = append(xmlResult, "</doStatement>")
		return xmlResult
	}
	eat(".", jackTokenizer, &xmlResult)
	jackTokenizer.Advance()
	utils.AppendIndent(&xmlResult, xmlToken(jackTokenizer.GetCurToken()))
	xmlResult = append(xmlResult, compileExpressionListInsideParentheses(jackTokenizer)...)
	eat(";", jackTokenizer, &xmlResult)
	xmlResult = append(xmlResult, "</doStatement>")
	return xmlResult
}

func CompileWhileStatement(jackTokenizer *la.JackTokenizer) []string {
	var xmlResult []string
	isExpectedToken("while", jackTokenizer.GetCurToken())
	xmlResult = append(xmlResult, "<whileStatement>")
	utils.AppendIndent(&xmlResult, xmlToken(jackTokenizer.GetCurToken()))
	eat("(", jackTokenizer, &xmlResult)
	jackTokenizer.Advance()
	utils.AppendIndent(&xmlResult, CompileExpression(jackTokenizer)...)

	eat(")", jackTokenizer, &xmlResult)
	eat("{", jackTokenizer, &xmlResult)
	for jackTokenizer.HasPeekToken() && !isOptionalExpectedToken("}", jackTokenizer.GetPeekToken()) {
		jackTokenizer.Advance()
		utils.AppendIndent(&xmlResult, CompileStatement(jackTokenizer.GetCurToken(), jackTokenizer)...)
	}
	eat("}", jackTokenizer, &xmlResult)
	xmlResult = append(xmlResult, "</whileStatement>")
	return xmlResult
}

func CompileVarDeclaration(jackTokenizer *la.JackTokenizer) []string {
	var xmlResult []string
	isExpectedToken("var", jackTokenizer.GetCurToken())
	xmlResult = append(xmlResult, "<varDec>")
	utils.AppendIndent(&xmlResult, xmlToken(jackTokenizer.GetCurToken()))
	jackTokenizer.Advance()
	utils.AppendIndent(&xmlResult, xmlToken(jackTokenizer.GetCurToken()))
	jackTokenizer.Advance()
	utils.AppendIndent(&xmlResult, xmlToken(jackTokenizer.GetCurToken()))
	eat(";", jackTokenizer, &xmlResult)
	xmlResult = append(xmlResult, "</varDec>")
	return xmlResult
}

func CompileParameterList(jackTokenizer *la.JackTokenizer) []string {
	var xmlResult []string
	xmlResult = append(xmlResult, "<parameterList>")
	peekToken := jackTokenizer.GetPeekToken()
	if peekToken.GetToken() == ")" {
		xmlResult = append(xmlResult, "</parameterList>")
		return xmlResult
	}
	utils.AppendIndent(&xmlResult, xmlToken(jackTokenizer.GetCurToken()))
	for jackTokenizer.HasPeekToken() && isOptionalExpectedToken(",", jackTokenizer.GetPeekToken()) {
		jackTokenizer.Advance()
		utils.AppendIndent(&xmlResult, xmlToken(jackTokenizer.GetCurToken()))
		jackTokenizer.Advance()
		utils.AppendIndent(&xmlResult, xmlToken(jackTokenizer.GetCurToken()))
	}
	xmlResult = append(xmlResult, "</parameterList>")
	return xmlResult
}

func compileStatements(jackTokenizer *la.JackTokenizer) []string {
	var xmlResult []string
	xmlResult = append(xmlResult, "<statements>")
	for jackTokenizer.HasPeekToken() && !isOptionalExpectedToken("}", jackTokenizer.GetPeekToken()) {
		jackTokenizer.Advance()
		xmlResult = append(xmlResult, CompileStatement(jackTokenizer.GetCurToken(), jackTokenizer)...)
	}
	xmlResult = append(xmlResult, "</statements>")
	return xmlResult
}

func compileSubroutineBody(jackTokenizer *la.JackTokenizer) []string {
	var xmlResult []string
	xmlResult = append(xmlResult, "<subroutineBody>")
	eat("{", jackTokenizer, &xmlResult)
	for jackTokenizer.HasPeekToken() && isOptionalExpectedToken("var", jackTokenizer.GetPeekToken()) {
		jackTokenizer.Advance()
		utils.AppendIndent(&xmlResult, CompileVarDeclaration(jackTokenizer)...)
	}
	utils.AppendIndent(&xmlResult, compileStatements(jackTokenizer)...)
	eat("}", jackTokenizer, &xmlResult)
	xmlResult = append(xmlResult, "</subroutineBody>")
	return xmlResult
}

func CompileSubroutine(jackTokenizer *la.JackTokenizer) []string {
	if !isOneOf(`constructor|function|method`, jackTokenizer.GetCurToken()) {
		panic(errors.New("Not a subroutine!"))
	}
	var xmlResult []string
	xmlResult = append(xmlResult, "<subroutineDec>")
	utils.AppendIndent(&xmlResult, xmlToken(jackTokenizer.GetCurToken()))
	jackTokenizer.Advance()
	utils.AppendIndent(&xmlResult, xmlToken(jackTokenizer.GetCurToken()))
	jackTokenizer.Advance()
	// if isOptionalExpectedToken("new", jackTokenizer.GetCurToken()) {
	// 	utils.AppendIndent(&xmlResult, XmlToken(jackTokenizer.GetCurToken()))
	// 	//jackTokenizer.Advance()
	// }
	utils.AppendIndent(&xmlResult, xmlToken(jackTokenizer.GetCurToken()))
	eat("(", jackTokenizer, &xmlResult)
	utils.AppendIndent(&xmlResult, CompileParameterList(jackTokenizer)...)
	eat(")", jackTokenizer, &xmlResult)
	utils.AppendIndent(&xmlResult, compileSubroutineBody(jackTokenizer)...)
	xmlResult = append(xmlResult, "</subroutineDec>")
	return xmlResult
}

func CompileClassVarDec(jackTokenizer *la.JackTokenizer) []string {
	if !isOneOf(`static|field`, jackTokenizer.GetCurToken()) {
		panic(errors.New("Not a subroutine!"))
	}
	var xmlResult []string
	xmlResult = append(xmlResult, "<classVarDec>")
	utils.AppendIndent(&xmlResult, xmlToken(jackTokenizer.GetCurToken()))
	jackTokenizer.Advance()
	utils.AppendIndent(&xmlResult, xmlToken(jackTokenizer.GetCurToken()))
	jackTokenizer.Advance()
	utils.AppendIndent(&xmlResult, xmlToken(jackTokenizer.GetCurToken()))
	eat(";", jackTokenizer, &xmlResult)
	xmlResult = append(xmlResult, "</classVarDec>")
	return xmlResult
}

func compileDeclaration(declaration la.Token, jackTokenize *la.JackTokenizer) ([]string, bool) {

	var xmlResult []string
	declarations := map[string]func(*la.JackTokenizer) []string{
		"field":       CompileClassVarDec,
		"static":      CompileClassVarDec,
		"function":    CompileSubroutine,
		"constructor": CompileSubroutine,
		"method":      CompileSubroutine,
	}

	compileDeclaration, valid := declarations[declaration.GetToken()]
	if valid {
		utils.AppendIndent(&xmlResult, compileDeclaration(jackTokenize)...)
	}
	return xmlResult, valid
}

func CompileClass(jackTokenizer *la.JackTokenizer) []string {
	isExpectedToken("class", jackTokenizer.GetCurToken())
	var xmlResult []string
	xmlResult = append(xmlResult, "<class>")
	utils.AppendIndent(&xmlResult, xmlToken(jackTokenizer.GetCurToken()))
	jackTokenizer.Advance()
	utils.AppendIndent(&xmlResult, xmlToken(jackTokenizer.GetCurToken()))
	eat("{", jackTokenizer, &xmlResult)
	totalCompiled := false
	for !totalCompiled {
		jackTokenizer.Advance()
		compiled, valid := compileDeclaration(jackTokenizer.GetCurToken(), jackTokenizer)
		if !valid {
			totalCompiled = true
		}
		xmlResult = append(xmlResult, compiled...)
	}
	utils.AppendIndent(&xmlResult, xmlToken(jackTokenizer.GetCurToken()))
	xmlResult = append(xmlResult, "</class>")
	return xmlResult
}
