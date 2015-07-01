package main

import "strings"

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

func Blockify(src string) GeneralNode {
    lines := strings.Split(src, "\n")
    idx := 0

    return parseGeneral(&idx, 0, lines)
}

func parseGeneral(idx *int, lvl int, lines []string) GeneralNode {
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

    return body
}

type StructureProvider struct {
    idx int
    body GeneralBody
}

func NewStructureProvider(body GeneralBody) *StructureProvider {
    return &StructureProvider{
        idx: 0,
        body: body,
    }
}

// return (line, exists)
func (s *StructureProvider) getLine() (string, bool) {
    if s.idx >= len(s.body.children) {
        return "", false
    }

    child := s.body.children[s.idx]
    switch child.(type) {
    case GeneralLine:
        s.idx++
        return child.(GeneralLine).line, true
    default:
        return "", false
    }
}

// return (block, exists)
func (s *StructureProvider) getBlock() (*StructureProvider, bool) {
    if s.idx >= len(s.body.children) {
        return nil, false
    }

    child := s.body.children[s.idx]
    switch child.(type) {
    case GeneralBody:
        s.idx++
        return NewStructureProvider(child.(GeneralBody)), true
    default:
        return nil, false
    }
}
