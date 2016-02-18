package parser

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestLexerValid(t *testing.T) {
	lexStr := "4172 \"test\" + - * / % () [] ! | . , : > >= < <= = == != -> identifier"

	lex := NewLexer(lexStr)
	go lex.Run()

	values := []string{
		"4172", "test", "+", "-", "*", "/", "%", "(", ")",
		"[", "]", "!", "|", ".", ",", ":", ">", ">=", "<",
		"<=", "=", "==", "!=", "->", "identifier"}
	types := []TokenType{
		TokenNumber, TokenString, TokenAdd, TokenSub, TokenMul,
		TokenDiv, TokenMod, TokenLParen, TokenRParen, TokenLBracket,
		TokenRBracket, TokenBang, TokenPipe, TokenDot, TokenComma,
		TokenColon, TokenGT, TokenGTE, TokenLT, TokenLTE,
		TokenAssign, TokenEQ, TokenNE, TokenArrow, TokenIdent}

	for i, val := range(values) {
		ty := types[i]

		tok := lex.NextToken()
		assert.Equal(t, val, tok.Value, "Incorrect token text")
		assert.Equal(t, ty, tok.Ty, "Incorrect token type")
	}

	tok := lex.NextToken()
	assert.Equal(t, TokenEOF, tok.Ty, "Expected to see EOF")
}

func TestLexerInvalid(t *testing.T) {
	lex := NewLexer("$")
	go lex.Run()

	tok := lex.NextToken()

	assert.Equal(t, TokenError, tok.Ty, "Expected error token")
}
