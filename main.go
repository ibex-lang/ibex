package main

import (
    "os"
    "io/ioutil"
    "log"

	"github.com/ibex-lang/ibex/parser"
)

func main() {
    args := os.Args[1:]

    for _, arg := range args {
        file, err := ioutil.ReadFile(arg)
        if err != nil {
            log.Print("Could not read file", arg)
            log.Print("Reason:", err)
        }

        compile(string(file), arg)
    }
}

func compile(src string, name string) {
    parser.InitExpressionParsing()

	body, err := parser.Blockify(src)
	if err != nil {
		log.Fatal(err)
	}
	structure := parser.NewStructure(body)
	ast, err := parser.Parse(structure)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%#v\n", ast)
}
