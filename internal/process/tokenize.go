package process

import (
	"fmt"
	"regexp"

	"cbeimers113/ilo/internal/config"
	"cbeimers113/ilo/internal/constant"
	"cbeimers113/ilo/internal/locale"
	"cbeimers113/ilo/internal/log"
)

var whitespaceRegex = regexp.MustCompile(`\s+`)

type Tokenizer struct {
	cfg   *config.Config
	input []rune
	pos   int  // Current position in input
	next  int  // Next position to read
	char  rune // Current character
	line  int  // Keep track of line number and position in line
	col   int
}

func NewTokenizer(cfg *config.Config, input string) *Tokenizer {
	t := &Tokenizer{
		cfg:   cfg,
		input: []rune(input),
		line:  1,
		col:   0,
	}

	t.readChar()
	return t
}

// NextToken reads the input char by char and produces the next token
func (t *Tokenizer) NextToken() Token {
	var tok Token
	t.skipWhitespace()

	if t.char == 0 {
		tok.Type = TokenEOF
	} else if isLetter(t.char) {
		literal := t.readIdentifier()
		return t.createToken(lookupIdent(literal), literal)
	} else if isDigit(t.char) {
		literal, tokType := t.readNumber()
		return t.createToken(tokType, literal)
	} else if isQuote(t.char) {
		literal, tokType := t.readQuote()
		return t.createToken(tokType, literal)
	} else if isComment(t.char) {
		literal := t.readComment()
		return t.createToken(TokenComment, literal)
	} else if constant.Operators.Contains(t.char) {
		op := string(t.char)
		t.readChar()
		return t.createToken(TokenOperator, op)
	} else {
		log.Fatal(fmt.Sprintf("%s %d: \"%s\"", t.cfg.Message(locale.ErrInvalidChar), t.line, string(t.char)))
	}

	t.readChar()
	return tok
}

func (t Tokenizer) createToken(tokType TokenType, literal string) Token {
	// If we're creating a string token, subtract the quotes from the column calculation
	colOffs := 0
	if tokType == TokenString || tokType == TokenFormString {
		colOffs = 2
	}

	return Token{
		Type:    tokType,
		Literal: literal,
		Line:    t.line,
		Col:     t.col - len(literal) - colOffs,
	}
}

// readChar reads the next character and advances position
func (t *Tokenizer) readChar() {
	if t.next >= len(t.input) {
		t.char = 0 // End of input
	} else {
		t.char = t.input[t.next]
	}

	t.pos = t.next
	t.next++
	t.col++
}

// readIdentifier reads the input into a buffer as long as the next char encountered is a letter.
// Returns the buffer when a non-letter is found
func (t *Tokenizer) readIdentifier() string {
	start := t.pos
	for isLetter(t.char) {
		t.readChar()
	}

	return string(t.input[start:t.pos])
}

// readNumber reads the input into a buffer as long as the next char encountered is a digit.
// Returns the buffer when a non-digit is found
func (t *Tokenizer) readNumber() (string, TokenType) {
	start := t.pos
	tokType := TokenInt

	// Determine if we're parsing an integer or float
	for isDigit(t.char) {
		t.readChar()

		// If we encounter a period, we should only keep reading if the following char is a digit
		if t.char == '.' && t.next < len(t.input) && isDigit(t.input[t.next]) {
			tokType = TokenFloat
			t.readChar()
		}
	}

	return string(t.input[start:t.pos]), tokType
}

// readQuote reads the input into a buffer as long as the next char encountered is not a matching quote.
// Returns the buffer when a matching quote is found. Panic if quote is unclosed
func (t *Tokenizer) readQuote() (string, TokenType) {
	quote := t.char
	line := t.line // Starting line of string
	tokType := TokenString

	if quote == '\'' {
		tokType = TokenFormString
	}

	t.readChar()
	start := t.pos

	for t.char != quote {
		t.readChar()

		// If we hit the end of the input without closing the quote, panic
		if t.char == 0 {
			log.Fatal(fmt.Sprintf("%s %d", t.cfg.Message(locale.ErrUnclosedQuote), line))
		}
	}

	end := t.pos
	t.readChar()
	return string(t.input[start:end]), tokType
}

// readComment reads the input into a buffer as long as the next char encountered is on the same line.
// Returns the buffer when the end of the line or EOF is reached
func (t *Tokenizer) readComment() string {
	start := t.pos
	for t.char != '\n' && t.char != 0 {
		t.readChar()
	}

	return string(t.input[start:t.pos])
}

// skipWhitespace reads over whitespace chars and advances the line count when a newline is found
func (t *Tokenizer) skipWhitespace() {
	for whitespaceRegex.MatchString(string(t.char)) {
		if t.char == '\n' {
			t.line++
			t.col = 0
		}
		t.readChar()
	}
}

// isLetter returns whether a given char is a letter
func isLetter(char rune) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

// isDigit returns whether a given char is a digit
func isDigit(char rune) bool {
	return '0' <= char && char <= '9'
}

// isQuote returns whether a given char is a quote
func isQuote(char rune) bool {
	return char == '"' || char == '\''
}

// isComment returns whether a given char is a comment initializer
func isComment(char rune) bool {
	return char == '~'
}

// lookupIdent checks if a given identifier is a keyword
func lookupIdent(ident string) TokenType {
	if constant.Keywords.Contains(ident) {
		if ident == "veras" || ident == "malveras" {
			return TokenBool
		}
		return TokenKeyword
	}

	return TokenIdent
}

// Tokenize takes raw source code and tokenizes it, returning a slice of the tokens encountered.
// Token reading will exit the program with a status code of 1 if a token-level syntax error occurs (unclosed quotes, invalid char, etc)
func Tokenize(cfg *config.Config, data string) []Token {
	tokenizer := NewTokenizer(cfg, data)
	tokens := make([]Token, 0)

	// Parse the raw source code into a list of tokens
	for token := tokenizer.NextToken(); token.Type != TokenEOF; token = tokenizer.NextToken() {
		if token.Literal == "" {
			continue
		}

		tokens = append(tokens, token)
	}

	return tokens
}
