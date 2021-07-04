package lexical_analyzer

type TokenType string

const (
	KEYWORD    TokenType = "keyword"
	SYMBOL               = "symbol"
	INTEGER              = "integerConstant"
	STRING               = "stringConstant"
	IDENTIFIER           = "identifier"
)
