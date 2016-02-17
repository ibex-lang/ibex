package parser

import (
    "strings"
    "fmt"
	"errors"

	"github.com/ibex-lang/ibex/core"
)

const indentWidth int = 4

func indentDepth(line string) (int, bool) {
    for i, c := range line {
        if c != ' ' {
            return i / indentWidth, i % indentWidth == 0
        }
    }
    return len(line) / indentWidth, len(line) % indentWidth == 0
}

// ugly algebraic types
type GeneralNode interface {
    isGeneral()
}

type GeneralBody struct {
    children []GeneralNode
}
func (g GeneralBody) isGeneral() {}

type GeneralLine struct {
    line string
}
func (g GeneralLine) isGeneral() {}

func Blockify(src string) (*GeneralBody, error) {
    lines := strings.Split(src, "\n")
    idx := 0

    return parseGeneral(&idx, 0, lines)
}

func parseGeneral(idx *int, lvl int, lines []string) (*GeneralBody, error) {
    body := GeneralBody{children: make([]GeneralNode, 0)}

    for *idx < len(lines) {
        line := lines[*idx]
        indent, valid := indentDepth(line)
		if !valid {
			return nil, errors.New("Invalid indentation")
		}

        if indent == lvl {
            child := GeneralLine{
                line: line[indent * indentWidth:],
            }
            body.children = append(body.children, child)
            *idx++
        } else if indent == lvl + 1 {
            child, err := parseGeneral(idx, lvl + 1, lines)
			if err != nil {
				return nil, err
			}
            body.children = append(body.children, child)
        } else {
            break
        }
    }

    return &body, nil
}

type ParseError struct {
    start int
    end int
    message string
}

func (e *ParseError) Error() string {
    return fmt.Sprintf(
        "Error at %d...%d: %s",
        e.start, e.end, e.message)
}

func ErrorAtToken(tok *Token, msg string) *ParseError {
    return &ParseError{
        start: tok.Start,
        end: tok.End,
        message: msg,
    }
}

type Structure struct {
    idx int
    body *GeneralBody
}

func NewStructure(body *GeneralBody) *Structure {
    return &Structure{
        idx: 0,
        body: body,
    }
}

// return (line, exists)
func (s *Structure) getLine() (*Lexer, bool) {
    if s.idx >= len(s.body.children) {
        return nil, false
    }

    child := s.body.children[s.idx]
    switch child.(type) {
    case GeneralLine:
        s.idx++
        lex := NewLexer(child.(GeneralLine).line)
        return lex, true
    default:
        return nil, false
    }
}

// return (block, exists)
func (s *Structure) getBlock() (*Structure, bool) {
    if s.idx >= len(s.body.children) {
        return nil, false
    }

    child := s.body.children[s.idx]
    switch child.(type) {
    case GeneralBody:
        s.idx++
        b := child.(GeneralBody)
        return NewStructure(&b), true
    default:
        return nil, false
    }
}

type ParsingContext int

const (
    ModuleContext = iota
    FunctionBodyContext
)
func ParseType(lex *Lexer) (core.IbexType, error) {
    return parseType(lex)
}

/*
func (s *Structure) parse(ctx ParsingContext) (ASTNode, error) {
    line, ok := s.getLine()
    for ok {
        line, ok = s.getLine()
    }

	return nil, nil
}*/

func parseType(lex *Lexer) (core.IbexType, error) {
    tok := lex.NextToken()
    switch tok.Ty {

    case TokenFunction:
        argType, err := parseType(lex)
        if err != nil {
            return nil, err
        }
        var retType core.IbexType = nil
        if lex.PeekToken().Ty == TokenArrow {
            lex.NextToken() // consume ->
            retType, err = parseType(lex)
            if err != nil {
                return nil, err
            }
        }
        return core.IbexFunctionType{argType, retType}, nil

    case TokenIdent:
        return parseIdentType(tok, lex)

    case TokenLParen:
        tok = lex.NextToken()

        namedTypes := make([]*core.IbexNamedTupleEntry, 0)
        normalTypes := make([]core.IbexType, 0)
        named := false
        if tok.Ty == TokenIdent {
            if lex.PeekToken().Ty == TokenColon {
                // named tuple
                lex.NextToken() // consume :
                tag := tok.Value
                named = true
                ty, err := parseType(lex)
                if err != nil {
                    return nil, err
                }
                entry := core.IbexNamedTupleEntry{tag, ty}
                namedTypes = append(namedTypes, &entry)
            } else {
                // normal tuple
                ty, err := parseIdentType(tok, lex)
                if err != nil {
                    return nil, err
                }
                normalTypes = append(normalTypes, ty)
            }
        } else {
            ty, err := parseType(lex)
            if err != nil {
                return nil, err
            }
            normalTypes = append(normalTypes, ty)
        }

        if named {
            for lex.PeekToken().Ty == TokenComma {
                lex.NextToken() // consume ,
                tok = lex.NextToken()
                if tok.Ty != TokenIdent {
                    return nil, ErrorAtToken(tok, "Expected identifier")
                }
                next := lex.NextToken()
                if next.Ty != TokenColon {
                    return nil, ErrorAtToken(tok, "Expected ':'")
                }

                ty, err := parseType(lex)
                if err != nil {
                    return nil, err
                }
                entry := core.IbexNamedTupleEntry{tok.Value, ty}
                namedTypes = append(namedTypes, &entry)
            }
            paren := lex.NextToken()
            if paren.Ty != TokenRParen {
                return nil, ErrorAtToken(paren, "Expected ')'")
            }
            return core.IbexNamedTupleType{namedTypes}, nil
        } else {
            for lex.PeekToken().Ty == TokenComma {
                lex.NextToken() // consume ,
                ty, err := parseType(lex)
                if err != nil {
                    return nil, err
                }
                normalTypes = append(normalTypes, ty)
            }
            paren := lex.NextToken()

            if paren.Ty != TokenRParen {
                return nil, ErrorAtToken(paren, "Expected ')'")
            }
            return core.IbexTupleType{normalTypes}, nil
        }

    case TokenLBracket:
        tok = lex.NextToken()
        if tok.Ty != TokenRBracket {
            return nil, ErrorAtToken(tok, "Expected ']'")
        }
        dims := 1
        for lex.PeekToken().Ty == TokenLBracket {
            lex.NextToken() // consume [
            tok = lex.NextToken()
            if tok.Ty != TokenRBracket {
               return nil, ErrorAtToken(tok, "Expected ']'")
            }
        }
        ty, err := parseType(lex)
        if err != nil {
            return nil, err
        }
        return core.IbexArrayType{ty, dims}, nil
    }

    return nil, ErrorAtToken(tok, "Unexpected token")
}

func parseIdentType(tok *Token, lex *Lexer) (core.IbexType, error) {
    if tok == nil {
        tok = lex.NextToken()
    }
    return core.IbexSimpleType{tok.Value}, nil
}

func Parse(s *Structure) (*ASTCompilationUnit, error) {
    return parse(s)
}

func parse(s *Structure) (*ASTCompilationUnit, error) {
    unit := &ASTCompilationUnit{
        Uses: make([]*ASTUseStmt, 0),
        Declarations: make([]ASTMemberDeclaration, 0),
    }

    lex, next := s.getLine()
    for next {
        go lex.Run()

        t := lex.NextToken()
        if t.Ty == TokenUse {
            use, err := parseUseStmt(lex)
            if err != nil {
                return nil, err
            }
            unit.Uses = append(unit.Uses, use)
        } else if t.Ty == TokenFunction {
            fn, err := parseFunction(lex)
            if err != nil {
                return nil, err
            }
            unit.Declarations = append(unit.Declarations, fn)
        } else if t.Ty == TokenTypeKW {
            decl, err := parseTypeDecl(lex)
            if err != nil {
                return nil, err
            }
            unit.Declarations = append(unit.Declarations, decl)
        }

        lex, next = s.getLine()
    }

    return unit, nil
}

func parseUseStmt(lex *Lexer) (*ASTUseStmt, error) {
    path := make([]string, 0)

    t := lex.NextToken()
    for t.Ty != TokenEOF {
        if t.Ty == TokenIdent {
            path = append(path, t.Value)
        } else if t.Ty != TokenModSep {
            return nil, ErrorAtToken(t, "Unexpected token")
        }

        t = lex.NextToken()
    }

    return &ASTUseStmt{path}, nil
}

func parseParameter(lex *Lexer) (*FunctionParameter, error) {
    tok := lex.NextToken()
    if tok.Ty != TokenIdent {
        return nil, ErrorAtToken(tok, "Expected identifier")
    }
    name := tok.Value
    tok = lex.NextToken()
    if tok.Ty != TokenColon {
        return nil, ErrorAtToken(tok, "Expected ':'")
    }
    ty, err := parseType(lex)
    if err != nil {
        return nil, err
    }
    return &FunctionParameter{name, ty}, nil
}

func parseFunction(lex *Lexer) (*ASTFunction, error) {
    ident := lex.NextToken()
    if ident.Ty != TokenIdent {
        return nil, ErrorAtToken(ident, "Expected identifier")
    }

    params := make([]*FunctionParameter, 0)
    peek := lex.PeekToken()
    if peek.Ty == TokenIdent {
        param, err := parseParameter(lex)
        if err != nil {
            return nil, err
        }
        params = append(params, param)
    } else if peek.Ty == TokenLParen {
        lex.NextToken() // consume (

        param, err := parseParameter(lex)
        if err != nil {
            return nil, err
        }
        params = append(params, param)

        for lex.PeekToken().Ty == TokenComma {
            lex.NextToken()
            param, err = parseParameter(lex)
            if err != nil {
                return nil, err
            }
            params = append(params, param)
        }

        paren := lex.NextToken()
        if paren.Ty != TokenRParen {
            return nil, ErrorAtToken(paren, "Expected ')'")
        }
    }

    var retType core.IbexType = nil
    tok := lex.NextToken()
    if tok.Ty == TokenArrow {
        ty, err := parseType(lex)
        if err != nil {
            return nil, err
        }
        retType = ty // doesn't compile without this!?
    }

    fn := ASTFunction{
        Name: ident.Value,
        Parameters: params,
        Return: retType,
        Body: nil,
    }
    return &fn, nil
}

func parseTypeDecl(lex *Lexer) (*ASTTypeDeclaration, error) {
    ident := lex.NextToken()
    if ident.Ty != TokenIdent {
        return nil, ErrorAtToken(ident, "Expected identifier")
    }

    tok := lex.NextToken()
    if tok.Ty != TokenAssign {
        return nil, ErrorAtToken(tok, "Expected '='")
    }

    ty, err := parseType(lex)
    if err != nil {
        return nil, err
    }

    decl := ASTTypeDeclaration{ident.Value, ty}
    return &decl, nil
}

