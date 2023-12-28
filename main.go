package main

import (
	"fmt"
	"monkey/repl"
	"os"
)

func main() {
	fmt.Println("Now you can use the Monkey programming language!")
	repl.Start(os.Stdin, os.Stdout)
}
