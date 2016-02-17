package main

import (
    "os"
    "io/ioutil"
    "fmt"

	"github.com/augustt198/ibex/parser"
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
    parser.InitExpressionParsing()

    lex := parser.NewLexer(src)
    go lex.Run()

    expr, err := parser.ParseExpression(lex)
    if err == nil {
        fmt.Printf("%#v\n", expr)
    } else {
        fmt.Println(err)
    }
}
