package main

type ASTNode interface {}

type ASTCompilationUnit struct {
    Uses []*ASTUseStmt
    Declarations []ASTMemberDeclaration
}

type ASTUseStmt struct {
    Path []string
}

type ASTMemberDeclaration interface {}

type ASTBody struct {
    Children []ASTNode
}

type ASTFunction struct {
    Name string
    Parameters []*FunctionParameter
    Return IbexType
    Body *ASTBody
}

type FunctionParameter struct {
    Name string
    Type IbexType
}

type Expression interface {
    ASTNode
}

type IdentExpr struct {
    Ident string
}

type NumberExpr struct {
    Number string
}

type NotExpr struct {
    Expr Expression
}

type NegateExpr struct {
    Expr Expression
}

type AddExpr struct {
    Left Expression
    Right Expression
}

type SubExpr struct {
    Left Expression
    Right Expression
}

type FunctionCallExpr struct {
    Input Expression
    Target Expression
}

type MulExpr struct {
    Left Expression
    Right Expression
}

type DivExpr struct {
    Left Expression
    Right Expression
}

type ModExpr struct {
    Left Expression
    Right Expression
}
