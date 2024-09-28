package process

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"

	"cbeimers113/ilo/internal/config"
	"cbeimers113/ilo/internal/util"
)

var (
	//go:embed test/test.ilo
	testDataGood string

	//go:embed test/test_invalid_char.ilo
	testDataInvalidChar string

	//go:embed test/test_unclosed_quote.ilo
	testDataUnclosedQuote string
)

func Test_NewTokenizer(t *testing.T) {
	cfg, err := config.TestConfig()
	assert.NoError(t, err)

	input := "hello, world!"
	tokenizer := NewTokenizer(cfg, input)
	assert.Equal(t, tokenizer.input, []rune(input))
	assert.Equal(t, 1, tokenizer.line)

	// col should advance to 1 when first char is read on initialization
	assert.Equal(t, 1, tokenizer.col)
	assert.Equal(t, 0, tokenizer.pos)
	assert.Equal(t, 1, tokenizer.next)
	assert.Equal(t, 'h', tokenizer.char)
}

func Test_NextToken(t *testing.T) {
	cfg, err := config.TestConfig()
	assert.NoError(t, err)

	// Create an input string with tokens of all types we expect
	// Also include different forms of whitespace to test that those are skipped
	input := "hello veras\tago ~comment\n2.5  7 \"str\"     'fstr' +"

	// These are the expected tokens from the above input
	want := []TokenType{
		TokenIdent,
		TokenBool,
		TokenKeyword,
		TokenComment,
		TokenFloat,
		TokenInt,
		TokenString,
		TokenFormString,
		TokenOperator,
	}

	// Record the tokens we encounter
	tokenizer := NewTokenizer(cfg, input)
	got := make([]TokenType, 0)
	for token := tokenizer.NextToken(); token.Type != TokenEOF; token = tokenizer.NextToken() {
		got = append(got, token.Type)
	}

	assert.Equal(t, want, got)
}

func Test_createToken(t *testing.T) {
	line := 5
	col := 7
	floatLiteral := "3.14159"
	stringLiteral := "hello, world!"

	tests := []struct {
		name    string
		tokType TokenType
		literal string
		want    Token
	}{
		{
			name:    "creates a non-string token",
			tokType: TokenFloat,
			literal: floatLiteral,
			want: Token{
				Type:    TokenFloat,
				Literal: floatLiteral,
				Line:    line,
				Col:     col - len(floatLiteral),
			},
		},
		{
			name:    "creates a string token with adjusted column position",
			tokType: TokenString,
			literal: stringLiteral,
			want: Token{
				Type:    TokenString,
				Literal: stringLiteral,
				Line:    line,
				Col:     col - len(stringLiteral) - 2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := &Tokenizer{
				line: line,
				col:  col,
			}

			got := tokenizer.createToken(tt.tokType, tt.literal)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_readChar(t *testing.T) {
	cfg, err := config.TestConfig()
	assert.NoError(t, err)

	input := "ch"
	tokenizer := NewTokenizer(cfg, input)
	assert.Equal(t, 'c', tokenizer.char)
	assert.Equal(t, 0, tokenizer.pos)
	assert.Equal(t, 1, tokenizer.next)
	assert.Equal(t, 1, tokenizer.col)

	tokenizer.readChar()
	assert.Equal(t, 'h', tokenizer.char)
	assert.Equal(t, 1, tokenizer.pos)
	assert.Equal(t, 2, tokenizer.next)
	assert.Equal(t, 2, tokenizer.col)

	tokenizer.readChar() // Reads EOF
	assert.Equal(t, rune(0), tokenizer.char)
	assert.Equal(t, 2, tokenizer.pos)
	assert.Equal(t, 3, tokenizer.next)
	assert.Equal(t, 3, tokenizer.col)
}

func Test_readIdentifier(t *testing.T) {
	cfg, err := config.TestConfig()
	assert.NoError(t, err)

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "reads an identifier",
			input: "hello",
			want:  "hello",
		},
		{
			name:  "doesn't read a number",
			input: "527",
		},
		{
			name:  "doesn't read a quote",
			input: "\"quote\"",
		},
		{
			name:  "doesn't read a comment",
			input: "~ comment",
		},
		{
			name:  "doesn't read whitespace",
			input: "\t   \n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := NewTokenizer(cfg, tt.input)
			got := tokenizer.readIdentifier()
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_readNumber(t *testing.T) {
	cfg, err := config.TestConfig()
	assert.NoError(t, err)

	tests := []struct {
		name      string
		input     string
		want      string
		wantFloat bool
	}{
		{
			name:  "reads an integer",
			input: "527",
			want:  "527",
		},
		{
			name:      "reads a float",
			input:     "3.14",
			want:      "3.14",
			wantFloat: true,
		},
		{
			name:  "reads an integer that is followed by a period",
			input: "54.",
			want:  "54",
		},
		{
			name:  "doesn't read an identifier",
			input: "hello",
		},
		{
			name:  "doesn't read a quote",
			input: "\"quote\"",
		},
		{
			name:  "doesn't read a comment",
			input: "~ comment",
		},
		{
			name:  "doesn't read whitespace",
			input: "\t   \n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// The default token type returned is int, since readNumber should only be called when encountering a digit
			wantType := TokenInt
			if tt.wantFloat {
				wantType = TokenFloat
			}

			tokenizer := NewTokenizer(cfg, tt.input)
			literal, tokType := tokenizer.readNumber()
			assert.Equal(t, tt.want, literal)
			assert.Equal(t, wantType, tokType)
		})
	}
}

func Test_readQuote(t *testing.T) {
	cfg, err := config.TestConfig()
	assert.NoError(t, err)

	tests := []struct {
		name        string
		input       string
		want        string
		wantFString bool
		wantExit    bool
	}{
		{
			name:  "reads a double quoted string",
			input: `"hello world" more tokens`,
			want:  "hello world",
		},
		{
			name:        "reads a single quoted string",
			input:       `'prints a {var}' more tokens`,
			want:        "prints a {var}",
			wantFString: true,
		},
		{
			name:     "errors out reading unclosed quote",
			input:    `"unclosed quote more tokens`,
			wantExit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := NewTokenizer(cfg, tt.input)

			if tt.wantExit {
				util.TestExits(t, func() { tokenizer.readQuote() }, "Test_readQuote")
			} else {
				wantType := TokenString
				if tt.wantFString {
					wantType = TokenFormString
				}

				got, tokType := tokenizer.readQuote()
				assert.Equal(t, tt.want, got)
				assert.Equal(t, wantType, tokType)
			}
		})
	}
}

func Test_readComment(t *testing.T) {
	cfg, err := config.TestConfig()
	assert.NoError(t, err)

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "reads comment until end of line",
			input: "~ this is a comment\nthis is not",
			want:  "~ this is a comment",
		},
		{
			name:  "reads comment at end of file",
			input: "~ this is the final comment",
			want:  "~ this is the final comment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := NewTokenizer(cfg, tt.input)
			got := tokenizer.readComment()
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_skipWhitespace(t *testing.T) {
	cfg, err := config.TestConfig()
	assert.NoError(t, err)

	tests := []struct {
		name     string
		input    string
		wantChar rune
		wantLine int
	}{
		{
			name:     "reads whitespace on same line",
			input:    "\t   h",
			wantChar: 'h',
			wantLine: 1,
		},
		{
			name:     "reads whitespace that advances line",
			input:    "\t\r\n    e",
			wantChar: 'e',
			wantLine: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := NewTokenizer(cfg, tt.input)
			tokenizer.skipWhitespace()
			assert.Equal(t, tt.wantChar, tokenizer.char)
			assert.Equal(t, tt.wantLine, tokenizer.line)
		})
	}
}

func Test_detectors(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantLetter  bool
		wantDigit   bool
		wantQuote   bool
		wantComment bool
	}{
		{
			name:       "detects letters",
			input:      "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_",
			wantLetter: true,
		},
		{
			name:      "detects digits",
			input:     "1234567890",
			wantDigit: true,
		},
		{
			name:      "detects quotes",
			input:     `"'`,
			wantQuote: true,
		},
		{
			name:        "detects comments",
			input:       "~",
			wantComment: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//lint:ignore S1029 we care more about rune position than byte position
			//lint:ignore SA6003 see above
			for _, char := range []rune(tt.input) {
				assert.Equal(t, tt.wantLetter, isLetter(char))
				assert.Equal(t, tt.wantDigit, isDigit(char))
				assert.Equal(t, tt.wantQuote, isQuote(char))
				assert.Equal(t, tt.wantComment, isComment(char))
			}
		})
	}
}

func Test_lookupIdent(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  TokenType
	}{
		{
			name:  "detects a boolean literal",
			input: "veras",
			want:  TokenBool,
		},
		{
			name:  "detects a keyword",
			input: "ago",
			want:  TokenKeyword,
		},
		{
			name:  "detects an identifier",
			input: "nomo",
			want:  TokenIdent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lookupIdent(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Tokenize(t *testing.T) {
	cfg, err := config.TestConfig()
	assert.NoError(t, err)

	tests := []struct {
		name       string
		data       string
		wantTokens []Token
		wantExit   bool
	}{
		{
			name: "tokenizes source code",
			data: testDataGood,
			wantTokens: []Token{
				{
					Type:    TokenComment,
					Literal: "~ Mia unua programo!",
					Line:    1,
					Col:     1,
				},
				{
					Type:    TokenKeyword,
					Literal: "jen",
					Line:    2,
					Col:     1,
				},
				{
					Type:    TokenKeyword,
					Literal: "ago",
					Line:    2,
					Col:     5,
				},
				{
					Type:    TokenIdent,
					Literal: "Komenci",
					Line:    2,
					Col:     9,
				},
				{
					Type:    TokenOperator,
					Literal: ":",
					Line:    2,
					Col:     16,
				},
				{
					Type:    TokenIdent,
					Literal: "Diru",
					Line:    3,
					Col:     5,
				},
				{
					Type:    TokenString,
					Literal: "Saluton mondo!",
					Line:    3,
					Col:     10,
				},
				{
					Type:    TokenOperator,
					Literal: ".",
					Line:    3,
					Col:     26,
				},
			},
		},
		{
			name:     "errors out on invalid char",
			data:     testDataInvalidChar,
			wantExit: true,
		},
		{
			name:     "errors out on unclosed quote",
			data:     testDataUnclosedQuote,
			wantExit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantExit {
				util.TestExits(t, func() { Tokenize(cfg, tt.data) }, "Test_Tokenize")
			} else {
				tokens := Tokenize(cfg, tt.data)
				assert.Equal(t, tt.wantTokens, tokens)
			}
		})
	}
}
