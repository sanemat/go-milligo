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

	fmt.Print("(module\n")
	fmt.Print("    (import \"wasi_unstable\" \"proc_exit\" (func $proc_exit (param i32)))\n")
	fmt.Print("    (memory 1)\n")
	fmt.Print("    (export \"memory\" (memory 0))\n")
	fmt.Print("    (func $main (export \"_start\")\n")
	fmt.Print("        (call $proc_exit\n")

	var outFirstNum = false
	var ns string
	var op rune
	for _, r := range os.Args[1] {
		if '0' <= r && r <= '9' {
			ns = ns + string(r)
		} else if r == '+' || r == '-' {
			if !outFirstNum {
				n, err := strconv.Atoi(ns)
				if err != nil {
					fmt.Printf("not integer. got=%s", ns)
				}
				fmt.Printf("            i32.const %d\n", n)
				ns = ""
				op = r
				outFirstNum = true
			} else {
				n, err := strconv.Atoi(ns)
				if err != nil {
					fmt.Printf("not integer. got=%s", ns)
				}
				fmt.Printf("            i32.const %d\n", n)
				ns = ""
				if op == '+' {
					fmt.Print("            i32.add\n")
				}
				if op == '-' {
					fmt.Print("            i32.sub\n")
				}
				op = r
			}

		}
	}
	if ns != "" {
		n, err := strconv.Atoi(ns)
		if err != nil {
			fmt.Printf("not integer. got=%s", ns)
		}
		fmt.Printf("            i32.const %d\n", n)
	}
	if op == '+' {
		fmt.Print("            i32.add\n")
	}
	if op == '-' {
		fmt.Print("            i32.sub\n")
	}

	fmt.Print("        )\n")
	fmt.Print("    )\n")
	fmt.Print(")\n")
}
