package lexical_analyzer

type TokenType string

const (
	KEYWORD    TokenType = "keyword"
	SYMBOL               = "symbol"
	INTEGER              = "integerConst"
	STRING               = "stringConst"
	IDENTIFIER           = "identifier"
)
