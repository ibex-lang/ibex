package main

import (
    "os"
    "io/ioutil"
    "fmt"
)

func main() {
    args := os.Args[1:]

    for _, arg := range args {
        file, err := ioutil.ReadFile(arg)
        if err != nil {
            fmt.Println("Could not read file", arg)
            fmt.Println("Reason:", err)
        }

        compile(string(file), arg)
    }
}

func compile(src string, name string) {
    InitExpressionParsing()

    lex := NewLexer(src)
    go lex.Run()

    expr, err := ParseExpression(lex)
    if err == nil {
        fmt.Printf("%#v\n", expr)
    } else {
        fmt.Println(err)
    }
}
