package main

import (
    "fmt"
    "unicode/utf8"
)

type TokenType int

const (
    TokenError TokenType = iota
    TokenEOF

    TokenIdent
    TokenNumber

    TokenFunction // fn
    TokenMatch    // match
    TokenUse      // use

    TokenAdd // +
    TokenSub // -
    TokenDiv // /
    TokenMul // *

    TokenBang   // !
    TokenPipe   // |
    TokenDot    // .
    TokenComma  // ,
    TokenColon  // :
    TokenModSep // ::

    TokenAssign // =
    TokenArrow  // ->

    TokenGT  // >
    TokenGTE // >=
    TokenLT  // <
    TokenLTE // <=
    TokenEQ  // ==
    TokenNE  // !=

    TokenLBrace // [
    TokenRBrace // ]
    TokenLParen // (
    TokenRParen // )
)

var keywords map[string]TokenType = map[string]TokenType{
    "fn": TokenFunction,
    "match": TokenMatch,
    "use": TokenUse,
}

type Token struct {
    Value string
    Ty    TokenType
    Start int
    End   int
}

type Lexer struct {
    src string
    start int
    pos int
    tokens chan *Token

    peekTok *Token // LL(1)
}

func NewLexer(src string) *Lexer {
    return &Lexer{
        src: src,
        start: 0,
        pos: 0,
        tokens: make(chan *Token),
    }
}

func (l *Lexer) NextToken() *Token {
    if l.peekTok != nil {
        tok := l.peekTok
        l.peekTok = nil
        return tok
    } else {
        return <-l.tokens
    }
}

func (l *Lexer) PeekToken() *Token {
    if l.peekTok != nil {
        return l.peekTok
    } else {
        l.peekTok = <-l.tokens
        return l.peekTok
    }
}

func (l *Lexer) peek() rune {
    if l.pos < len(l.src) {
        chr, _ := utf8.DecodeRuneInString(l.src[l.pos:])
        return chr
    }
    return utf8.RuneError 
}

func (l *Lexer) read() rune {
    if l.pos < len(l.src) {
        chr, size := utf8.DecodeRuneInString(l.src[l.pos:])
        l.pos += size
        return chr
    }
    return utf8.RuneError
}

func (l *Lexer) accept(chr rune) bool {
    if l.peek() == chr {
        l.read()
        return true
    }
    return false
}

func (l *Lexer) acceptAny(runes string) bool {
    chr := l.peek()
    for _, c := range runes {
        if chr == c {
            l.read()
            return true
        }
    }
    return false
}

func (l *Lexer) emitError(msg string) {
    l.tokens <- &Token{
        Value: msg,
        Ty: TokenError,
        Start: l.start,
        End: l.pos,
    }
    l.start = l.pos
}

func (l *Lexer) emitToken(ty TokenType) {
    l.tokens <- &Token{
        Value: l.src[l.start:l.pos],
        Ty: ty,
        Start: l.start,
        End: l.pos,
    }
    l.start = l.pos
}

func (l *Lexer) Run() {
    for l.pos < len(l.src) {
        if !l.getToken() {
            close(l.tokens)
            return
        }
    }
    l.emitToken(TokenEOF)
    close(l.tokens)
}

// ret = success?
func (l *Lexer) getToken() bool {
    switch chr := l.read(); chr {
    case '+': l.emitToken(TokenAdd)
    case '-':
        if l.accept('>') {
            l.emitToken(TokenArrow)
        } else {
            l.emitToken(TokenSub)
        }
    case '/': l.emitToken(TokenDiv)
    case '*': l.emitToken(TokenMul)

    case '!':
        if l.accept('=') {
            l.emitToken(TokenNE)
        } else {
            l.emitToken(TokenBang)
        }
    case '|': l.emitToken(TokenPipe)
    case '.': l.emitToken(TokenDot)
    case ',': l.emitToken(TokenComma)
    case ':':
        if l.accept(':') {
            l.emitToken(TokenModSep)
        } else {
            l.emitToken(TokenColon)
        }

    case '>':
        if l.accept('=') {
            l.emitToken(TokenGTE)
        } else {
            l.emitToken(TokenGT)
        }
    case '<':
        if l.accept('=') {
            l.emitToken(TokenLTE)
        } else {
            l.emitToken(TokenLT)
        }
    case '=':
        if l.accept('=') {
            l.emitToken(TokenEQ)
        } else {
            l.emitToken(TokenAssign)
        }

    case '[': l.emitToken(TokenLBrace)
    case ']': l.emitToken(TokenRBrace)
    case '(': l.emitToken(TokenLParen)
    case ')': l.emitToken(TokenRParen)

    case ' ', '\n', '\r':
        l.start++
        break

    default:
        if IsIdentStart(chr) {
            l.readIdent()
        } else if IsDigit(chr) {
            l.readNumber()            
        } else {
            err := fmt.Sprintf("Unexpected character: '%c'", chr)
            l.emitError(err)
            return false
        }
    }
    return true
}

func (l *Lexer) readIdent() {
    for IsIdentChar(l.peek()) {
        l.read()
    }
    ident := l.src[l.start:l.pos]

    keyword, exist := keywords[ident]
    if exist {
        l.emitToken(keyword)
    } else {
        l.emitToken(TokenIdent)
    }
}

func (l *Lexer) readNumber() {
    for IsDigit(l.peek()) {
        l.read()
    }
    l.emitToken(TokenNumber)
}
