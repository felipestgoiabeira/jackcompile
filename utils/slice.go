package utils

import "jackcompile/lexical_analysis"

func ContainsTokenType(s []lexical_analysis.TokenType, e lexical_analysis.TokenType) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
