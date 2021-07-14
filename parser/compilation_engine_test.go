package parser

import (
	"fmt"
	la "jackcompile/lexical_analysis"
	"testing"
)

func TestCompileTerm(t *testing.T) {
	token := la.NewToken("test")
	result := CompileTerm(*token)
	for _, r := range result {
		fmt.Println(r)
	}
}

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
