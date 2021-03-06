package parser

const (
    _ = iota // ignore 0
    AssignmentPrecedence     // =
    FunctionCallPrecedence   // ->
    AdditivePrecedence       // +, -
    MultiplicativePrecedence // *, /, %
    PostfixPrecedence
)

var prefixParsers map[TokenType]PrefixParser
var infixParsers map[TokenType]InfixParser
var postfixParsers map[TokenType]PostfixParser

// bypass initialization loop
func InitExpressionParsing() {
    prefixParsers = map[TokenType]PrefixParser{
        TokenIdent:     ParseIdent,
        TokenString:    ParseString,
        TokenNumber:    ParseNumber,
        TokenBang:      ParseUnaryPrefix,
        TokenSub:       ParseUnaryPrefix,
        TokenLParen:    ParseGrouping,
    }

    additive := InfixParser{ParseAdditive, AdditivePrecedence}
    multiplicative := InfixParser{ParseMultiplicative, MultiplicativePrecedence}

    infixParsers = map[TokenType]InfixParser{
        TokenAdd:   additive,
        TokenSub:   additive,
        TokenMul:   multiplicative,
        TokenDiv:   multiplicative,
        TokenMod:   multiplicative,
        TokenArrow: InfixParser{ParseFunctionCall, FunctionCallPrecedence},
    }

    postfixParsers = map[TokenType]PostfixParser{
        TokenBang:   ParseUnsafeAccess,
        TokenLBracket: ParseArrayAccess,
    }
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

    _, next := postfixParsers[lex.PeekToken().Ty]
    for next {
        tok := lex.NextToken()
        left, err = postfixParsers[tok.Ty](left, lex, tok)
        if err != nil {
            return nil, err
        }

        _, next = postfixParsers[lex.PeekToken().Ty]
    }

    for precedence < nextPrecedence(lex) {
        tok = lex.NextToken()

        parser := infixParsers[tok.Ty]
        expr, err := parser.Parser(left, lex, tok)
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
        return parser.Precedence
    } else {
        return 0
    }
}

type PrefixParser func (*Lexer, *Token) (Expression, error)

func ParseIdent(lex *Lexer, tok *Token) (Expression, error) {
    return IdentExpr{tok.Value}, nil
}

func ParseString(lex *Lexer, tok *Token) (Expression, error) {
    return StringExpr{tok.Value}, nil
}

func ParseNumber(lex *Lexer, tok *Token) (Expression, error) {
    return NumberExpr{tok.Value}, nil
}

func ParseUnaryPrefix(lex *Lexer, tok *Token) (Expression, error) {
    expr, err := ParseExpression(lex)
    if err != nil {
        return nil, err
    }

    if tok.Ty == TokenBang {
        return NotExpr{expr}, nil
    } else if tok.Ty == TokenSub {
        return NegateExpr{expr}, nil
    } else {
        return nil, ErrorAtToken(tok, "Unexpected token")
    }
}

// parses a grouping/tuple/named tuple
func ParseGrouping(lex *Lexer, tok *Token) (Expression, error) {
    expr, err := ParseExpression(lex)
    if err != nil {
        return nil, err
    }

    peek := lex.PeekToken()
    if peek.Ty == TokenRParen {
		lex.NextToken()
        return expr, nil
    } else if peek.Ty == TokenComma {
        return parseTupleLiteral(expr, lex)
    } else if peek.Ty == TokenColon {
        return parseNamedTupleLiteral(expr, lex)
    } else {
        return nil, ErrorAtToken(peek, "Expected ')'")
    }
}

func parseTupleLiteral(first Expression, lex *Lexer) (Expression, error) {
    elems := []Expression{first}
    tok := lex.NextToken()
    for tok.Ty == TokenComma {
        expr, err := ParseExpression(lex)
        if err != nil {
            return nil, err
        }
        elems = append(elems, expr)
        tok = lex.NextToken()
    }

    if tok.Ty != TokenRParen {
        return nil, ErrorAtToken(tok, "Expected ')'")
    }

    return TupleExpr{elems}, nil
}

func parseNamedTupleLiteral(firstTag Expression,
    lex *Lexer) (Expression, error) {

    lex.NextToken()
    var tag string
    // TODO make better
    switch firstTag.(type) {
    case IdentExpr:
        tag = firstTag.(IdentExpr).Ident
    default:
        // PeekToken() should = colon
        return nil, ErrorAtToken(lex.PeekToken(), "Expected identifier preceding")
    }

    expr, err := ParseExpression(lex)
    if err != nil {
        return nil, err
    }

    elems := []*NamedTupleEntry{&NamedTupleEntry{tag, expr}}

    tok := lex.NextToken()
    for tok.Ty == TokenComma {
        tok = lex.NextToken()
        if tok.Ty != TokenIdent {
            return nil, ErrorAtToken(tok, "Expected identifier")
        }
        tag = tok.Value
        tok = lex.NextToken()
        if tok.Ty != TokenColon {
            return nil, ErrorAtToken(tok, "Expected ':'")
        }
        expr, err := ParseExpression(lex)
        if err != nil {
            return nil, err
        }
        elems = append(elems, &NamedTupleEntry{tag, expr})
        tok = lex.NextToken()
    }
    if tok.Ty != TokenRParen {
        return nil, ErrorAtToken(tok, "Expected ')'")
    }

    return NamedTupleExpr{elems}, nil
}

type InfixParser struct {
    Parser func(Expression, *Lexer, *Token) (Expression, error)
    Precedence int
}

type PostfixParser func(Expression, *Lexer, *Token) (Expression, error)

func ParseAdditive(left Expression, lex *Lexer,
    tok *Token) (Expression, error) {

    right, err := ParseExpressionP(AdditivePrecedence, lex)
    if err != nil {
        return nil, err
    }

    if tok.Ty == TokenAdd {
        return AddExpr{left, right}, nil
    } else if tok.Ty == TokenSub {
        return SubExpr{left, right}, nil
    } else {
        return nil, ErrorAtToken(tok, "Unexpected token")
    }
}

func ParseFunctionCall(left Expression, lex *Lexer,
    tok *Token) (Expression, error) {

    right, err := ParseExpressionP(FunctionCallPrecedence, lex)
    if err != nil {
        return nil, err
    }

    return FunctionCallExpr{left, right}, nil
}

func ParseMultiplicative(left Expression, lex *Lexer,
    tok *Token) (Expression, error) {

    right, err := ParseExpressionP(MultiplicativePrecedence, lex)
    if err != nil {
        return nil, err
    }

    if tok.Ty == TokenMul {
        return MulExpr{left, right}, nil
    } else if tok.Ty == TokenDiv {
        return DivExpr{left, right}, nil
    } else if tok.Ty == TokenMod {
        return ModExpr{left, right}, nil
    } else {
        return nil, ErrorAtToken(tok, "Unexpected token")
    }
}

func ParseUnsafeAccess(left Expression, lex *Lexer,
    tok *Token) (Expression, error) {

    return UnsafeAccessExpr{left}, nil
}

func ParseArrayAccess(left Expression, lex *Lexer,
    tok *Token) (Expression, error) {

    idx, err := ParseExpression(lex)
    if err != nil {
        return nil, err
    }

    tok = lex.NextToken()
    if tok.Ty != TokenRBracket {
        return nil, ErrorAtToken(tok, "Expected ']'")
    }

    return ArrayAccessExpr{left, idx}, nil
}
