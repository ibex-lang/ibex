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
    Body ASTBody
}
