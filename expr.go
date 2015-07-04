package main

var prefixParsers map[TokenType]PrefixParser
var infixParsers map[TokenType]InfixParser

// bypass initialization loop
func InitExpressionParsing() {
    prefixParsers = map[TokenType]PrefixParser{
        TokenIdent:     ParseIdent,
        TokenNumber:    ParseNumber,
        TokenBang:      ParseUnaryPrefix,
        TokenSub:       ParseUnaryPrefix,
    }

    infixParsers = map[TokenType]InfixParser{}
}

func ParseExpression(lex *Lexer) (Expression, error) {
    return ParseExpressionP(0, lex)
}

func ParseExpressionP(precedence int, lex *Lexer) (Expression, error) {
    tok := lex.NextToken()

    prefixParser, ok := prefixParsers[tok.Ty]
    if !ok {
        return nil, ErrorAtToken(tok, "Unexpected token")
    }

    left, err := prefixParser(lex, tok)
    if err != nil {
        return nil, err
    }

    for precedence < nextPrecedence(lex) {
        tok = lex.NextToken()

        parser := infixParsers[tok.Ty]
        expr, err := parser.Parse(left, lex, tok)
        if err != nil {
            return nil, err
        }
        left = expr
    }

    return left, nil
}

func nextPrecedence(lex *Lexer) int {
    parser, ok := infixParsers[lex.PeekToken().Ty]
    if ok {
        return parser.Precedence()
    } else {
        return 0
    }
}

type InfixParser interface {
    Parse(left Expression, lex *Lexer, tok *Token) (Expression, error)
    Precedence() int
}

type PrefixParser func (*Lexer, *Token) (Expression, error)


func ParseIdent(lex *Lexer, tok *Token) (Expression, error) {
    return IdentExpr{tok.Value}, nil
}

func ParseNumber(lex *Lexer, tok *Token) (Expression, error) {
    return NumberExpr{tok.Value}, nil
}

func ParseUnaryPrefix(lex *Lexer, tok *Token) (Expression, error) {
    expr, err := ParseExpression(lex)
    if err != nil {
        return nil, err
    }
    switch tok.Ty {
    case TokenBang:
        return NotExpr{expr}, nil

    case TokenSub:
        return NegateExpr{expr}, nil
    }

    return nil, ErrorAtToken(tok, "Unexpected token")
}

func ParseGrouping(lex *Lexer, tok *Token) (Expression, error) {
    expr, err := ParseExpression(lex)
    if err != nil {
        return nil, err
    }
    paren := lex.NextToken()
    if paren.Ty != TokenRParen {
        return nil, ErrorAtToken(paren, "Expected ')'")
    }

    return expr, nil
}
