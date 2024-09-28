package process

import (
	"fmt"

	"cbeimers113/ilo/internal/constant"
)

type TokenType string

const (
	TokenEOF        TokenType = "EOF"
	TokenKeyword    TokenType = "KEYWORD"
	TokenOperator   TokenType = "OPERATOR"
	TokenIdent      TokenType = "IDENTIFIER"
	TokenBool       TokenType = "BOOLEAN"
	TokenInt        TokenType = "INTEGER"
	TokenFloat      TokenType = "FLOAT"
	TokenString     TokenType = "STRING"
	TokenFormString TokenType = "FORMSTRING"
	TokenComment    TokenType = "COMMENT"
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Col     int
}

func (t Token) String() string {
	return fmt.Sprintf(
		"{%s%s %s%s%s Line: %s%d%s, Col: %s%d%s}",
		constant.ColBlue,
		t.Type,
		constant.ColGreen,
		t.Literal,
		constant.ColReset,
		constant.ColYellow,
		t.Line,
		constant.ColReset,
		constant.ColYellow,
		t.Col,
		constant.ColReset,
	)
}
