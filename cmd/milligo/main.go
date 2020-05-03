package main

import "fmt"

func main() {
	fmt.Print("(module\n")
	fmt.Print("    (import \"wasi_unstable\" \"proc_exit\" (func $proc_exit (param i32)))\n")
	fmt.Print("    (memory 1)\n")
	fmt.Print("    (export \"memory\" (memory 0))\n")
	fmt.Print("    (func $main (export \"_start\")\n")
	fmt.Print("        (call $proc_exit\n")
	fmt.Print("            (i32.const 42)\n")
	fmt.Print("        )\n")
	fmt.Print("    )\n")
	fmt.Print(")\n")
}
