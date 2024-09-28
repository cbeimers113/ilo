package process

import (
	"cbeimers113/ilo/internal/config"
	"cbeimers113/ilo/internal/locale"
	"cbeimers113/ilo/internal/log"
	"fmt"
)

// isOperator checks if a given token is an operator
func isOperator(token Token) bool {
	return token.Type == TokenOperator
}

// parsePrimary parses the atomic elements of the source code (literals & variables)
func parsePrimary(cfg *config.Config, tokens []Token, pos *int) ASTNode {
	token := tokens[*pos]
	*pos++ // Move to the next token

	switch token.Type {
	case TokenInt, TokenFloat, TokenString, TokenBool:
		return &Literal{Value: token}
	case TokenIdent:
		return &Variable{Name: token}
	}

	log.Fatal(fmt.Sprintf("%s %d: \"%s\"", cfg.Message(locale.ErrInvalidToken), token.Line, string(token.Literal)))
	return nil
}

// parseExpression parses a binary expression
func parseExpression(cfg *config.Config, tokens []Token, pos *int) ASTNode {
	left := parsePrimary(cfg, tokens, pos)

	for isOperator(tokens[*pos]) {
		op := tokens[*pos]
		*pos++ // Move to the next token
		right := parsePrimary(cfg, tokens, pos)
		left = &BinaryExpression{Left: left, Op: op, Right: right}
	}

	return left
}

// parseAssignment parses a variable assignment
func parseAssignment(cfg *config.Config, tokens []Token, pos *int) *Assignment {
	// Expect variable name
	variable := parsePrimary(cfg, tokens, pos).(*Variable)

	// Expect '=' operator
	if tokens[*pos].Type != TokenOperator || tokens[*pos].Literal != "=" {
		panic("Expected '=' in assignment")
	}
	*pos++ // Move past '='

	// Parse the right-hand side expression
	value := parseExpression(cfg, tokens, pos)

	return &Assignment{Variable: variable, Value: value}
}
