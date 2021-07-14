package parser

import (
	la "jackcompile/lexical_analysis"
)

type JackCompile struct {
	result []string
}

func NewJackCompile(file string) *JackCompile {
	j := new(JackCompile)
	jackTokenizer := la.NewJackTokenizer(file)
	j.result = CompileClass(jackTokenizer)
	return j
}

func (j *JackCompile) GetResult() []string {
	return j.result
}
