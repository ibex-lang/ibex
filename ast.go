package main

type ASTNode interface {}

type ASTCompilationUnit struct {
    Uses []*ASTUseStmt
    Declarations []ASTMemberDeclaration
}

type ASTMemberDeclaration interface {}

type ASTUseStmt struct {
    Path []string
}

type ASTBody struct {
    Children []ASTNode
}

type ASTFunction struct {
    Name string
    Body ASTBody
}
