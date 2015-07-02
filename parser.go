package main

import (
    "strings"
    "fmt"
)

const indentWidth int = 4

func indentDepth(line string) int {
    for i, c := range line {
        if c != ' ' {
            return i / indentWidth
        }
    }
    return len(line) / indentWidth
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

func Blockify(src string) *GeneralBody {
    lines := strings.Split(src, "\n")
    idx := 0

    return parseGeneral(&idx, 0, lines)
}

func parseGeneral(idx *int, lvl int, lines []string) *GeneralBody {
    body := GeneralBody{children: make([]GeneralNode, 0)}

    for *idx < len(lines) {
        line := lines[*idx]
        indent := indentDepth(line)
        if indent == lvl {
            child := GeneralLine{
                line: line[indent * indentWidth:],
            }
            body.children = append(body.children, child)
            *idx++
        } else if indent == lvl + 1 {
            child := parseGeneral(idx, lvl + 1, lines)
            body.children = append(body.children, child)
        } else {
            break
        }
    }

    return &body
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
