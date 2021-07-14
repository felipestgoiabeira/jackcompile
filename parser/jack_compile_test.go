package parser

import (
	"bytes"
	"fmt"
	"jackcompile/utils"
	"testing"
)

func TestJackCompile(t *testing.T) {
	jackcompile := NewJackCompile("../resources/tests/parser/LessSquare.jack")

	result := jackcompile.GetResult()
	var resultBuffer bytes.Buffer
	for _, r := range result {
		resultBuffer.WriteString(r + "\n")
		fmt.Println(r)
	}

	utils.WriteResultToFile(resultBuffer, "LessSquare.xml")
}
