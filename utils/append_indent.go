package utils

func AppendIndent(result *[]string, toAppend ...string) {
	for _, term := range toAppend {
		*result = append(*result, "  "+term)
	}
}
