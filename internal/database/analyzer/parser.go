package analyzer

import (
	"strings"
)

type TokenType int

const (
	TokenKeyword TokenType = iota
	TokenIdentifier
	TokenLiteral
	TokenOperator
	TokenSymbol
)

type Token struct {
	Type  TokenType
	Value string
}

// Tokenize performs basic SQL lexical analysis
func Tokenize(sql string) []Token {
	var tokens []Token
	words := strings.Fields(sql)

	keywords := map[string]bool{
		"SELECT": true, "INSERT": true, "UPDATE": true, "DELETE": true,
		"FROM": true, "WHERE": true, "JOIN": true, "DROP": true,
		"ALTER": true, "CREATE": true, "GRANT": true, "REVOKE": true,
		"UNION": true, "ALL": true, "EXEC": true, "EXECUTE": true,
	}

	for _, w := range words {
		clean := strings.Trim(w, ";(),")
		upper := strings.ToUpper(clean)
		if keywords[upper] {
			tokens = append(tokens, Token{Type: TokenKeyword, Value: upper})
		} else if strings.HasPrefix(clean, "'") || strings.HasPrefix(clean, "\"") {
			tokens = append(tokens, Token{Type: TokenLiteral, Value: clean})
		} else {
			tokens = append(tokens, Token{Type: TokenIdentifier, Value: clean})
		}
	}
	return tokens
}
