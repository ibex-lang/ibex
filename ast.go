package main

type ASTNode interface {}

type ASTBody struct {
    Children []ASTNode
}

type ASTFunction struct {
    Name string
    Body ASTBody
}
