package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	la "jackcompile/lexical_analysis"
	"jackcompile/utils"
	"testing"
)

func TestCompileExpression(t *testing.T) {
	jackTokenizer := la.NewJackTokenizer("../resources/tests/parser/expression.jack")
	result := CompileExpression(jackTokenizer)
	for _, r := range result {
		fmt.Println(r)
	}
}

func TestIfStatement(t *testing.T) {
	jackTokenizer := la.NewJackTokenizer("../resources/tests/parser/ifStatement.jack")
	result := CompileIfStatement(jackTokenizer)
	for _, r := range result {
		fmt.Println(r)
	}
}

func TestLetStatement(t *testing.T) {
	jackTokenizer := la.NewJackTokenizer("../resources/tests/parser/letStatement.jack")
	result := CompileLetStatement(jackTokenizer)
	for _, r := range result {
		fmt.Println(r)
	}
}

func TestReturnStatement(t *testing.T) {
	jackTokenizer := la.NewJackTokenizer("../resources/tests/parser/returnStatement.jack")
	result := CompileReturnStatement(jackTokenizer)
	for _, r := range result {
		fmt.Println(r)
	}
}

func TestCompileExpressionList(t *testing.T) {
	jackTokenizer := la.NewJackTokenizer("../resources/tests/parser/expressionList.jack")
	result := CompileExpressionList(jackTokenizer)
	for _, r := range result {
		fmt.Println(r)
	}
}

func TestCompileDoStatement(t *testing.T) {
	jackTokenizer := la.NewJackTokenizer("../resources/tests/parser/doStatement.jack")
	result := CompileDoStatement(jackTokenizer)
	for _, r := range result {
		fmt.Println(r)
	}
}

func TestCompileWhileStatement(t *testing.T) {
	jackTokenizer := la.NewJackTokenizer("../resources/tests/parser/whileStatement.jack")
	result := CompileWhileStatement(jackTokenizer)
	for _, r := range result {
		fmt.Println(r)
	}
}
func TestCompileVarDeclaration(t *testing.T) {
	jackTokenizer := la.NewJackTokenizer("../resources/tests/parser/varDec.jack")
	result := CompileVarDeclaration(jackTokenizer)
	for _, r := range result {
		fmt.Println(r)
	}
}
func TestCompileParameterList(t *testing.T) {
	jackTokenizer := la.NewJackTokenizer("../resources/tests/parser/parameterList.jack")
	result := CompileParameterList(jackTokenizer)
	for _, r := range result {
		fmt.Println(r)
	}
}

func TestCompileSubroutine(t *testing.T) {
	jackTokenizer := la.NewJackTokenizer("../resources/tests/parser/subroutine.jack")
	result := CompileSubroutine(jackTokenizer)
	for _, r := range result {
		fmt.Println(r)
	}
}

func TestCompileClassVarDeclaration(t *testing.T) {
	jackTokenizer := la.NewJackTokenizer("../resources/tests/parser/classVarDec.jack")
	result := CompileClassVarDec(jackTokenizer)
	for _, r := range result {
		fmt.Println(r)
	}
}

func TestCompileClass(t *testing.T) {
	jackTokenizer := la.NewJackTokenizer("../resources/tests/parser/Square.jack")
	result := CompileClass(jackTokenizer)
	var resultBuffer bytes.Buffer
	for _, r := range result {
		resultBuffer.WriteString(r + "\n")
		fmt.Println(r)
	}

	content, err := ioutil.ReadFile("../resources/xml/Square.xml")

	if err != nil {
		panic(err)
	}

	expectedResult := string(content)
	resultString := resultBuffer.String()
	if expectedResult != resultString {
		panic(errors.New("The result isn't valid!"))
	}

	utils.WriteResultToFile(resultBuffer, "Square.xml")
}
