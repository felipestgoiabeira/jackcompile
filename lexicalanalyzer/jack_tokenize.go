package lexicalanalyzer

type JackTokenize interface {
	hasMoreToken() bool
	tokenType() int
	keyWord() string
	symbol() string
	identifier() string
	intVal() int
	stringVal() string
}
