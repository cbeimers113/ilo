package abstract

const (
	TokenKeyword = iota
	TokenIdentifier
	TokenOperator
	TokenLiteral
	TokenPunctuation
)

type Token struct {
	ID string
	Type int
}