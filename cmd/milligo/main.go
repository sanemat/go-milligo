package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("args is not 2. got=%d", len(os.Args))
		os.Exit(1)
	}
	rv, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("args[1] is not integer. got=%s", os.Args[1])
		os.Exit(1)
	}

	fmt.Print("(module\n")
	fmt.Print("    (import \"wasi_unstable\" \"proc_exit\" (func $proc_exit (param i32)))\n")
	fmt.Print("    (memory 1)\n")
	fmt.Print("    (export \"memory\" (memory 0))\n")
	fmt.Print("    (func $main (export \"_start\")\n")
	fmt.Print("        (call $proc_exit\n")
	fmt.Printf("            (i32.const %d)\n", rv)
	fmt.Print("        )\n")
	fmt.Print("    )\n")
	fmt.Print(")\n")
}
