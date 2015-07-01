package main

type TokenType int

const (
    TokenError TokenType = iota

    TokenEOF
    TokenIdent
    TokenFunction
    TokenMatch
    TokenUse
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
