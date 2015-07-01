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
    block := Blockify(src)
    fmt.Printf("%s:\n", name)
    fmt.Println(block)
}
