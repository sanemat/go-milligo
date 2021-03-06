package codegen

import (
	"fmt"
	"github.com/sanemat/go-milligo/astnode"
)

//
// Code generator
//
func gen(node *astnode.Astnode) {
	switch node.Kind {
	case astnode.NUM:
		fmt.Printf("        i32.const %d\n", node.Val)
		return
	case astnode.RETURN:
		gen(node.LHS)
		fmt.Print("        return\n")
		return
	}

	gen(node.LHS)
	gen(node.RHS)

	switch node.Kind {
	case astnode.ADD:
		fmt.Print("        i32.add\n")
	case astnode.SUB:
		fmt.Print("        i32.sub\n")
	case astnode.MUL:
		fmt.Print("        i32.mul\n")
	case astnode.DIV:
		fmt.Print("        i32.div_s\n")
	case astnode.EQ:
		fmt.Print("        i32.eq\n")
	case astnode.NE:
		fmt.Print("        i32.ne\n")
	case astnode.LT:
		fmt.Print("        i32.lt_s\n")
	case astnode.LE:
		fmt.Print("        i32.le_s\n")
	}
}

// Codegen generates code
func Codegen(codes []*astnode.Astnode) {
	fmt.Print("(module\n")
	fmt.Print("    (import \"wasi_unstable\" \"proc_exit\" (func $proc_exit (param i32)))\n")
	fmt.Print("    (memory 1)\n")
	fmt.Print("    (export \"memory\" (memory 0))\n")
	fmt.Print("    (func $__ (result i32)\n")
	fmt.Print("        (local $tmp i32)\n")

	for _, code := range codes {
		// Traverse the AST to emit assembly.
		gen(code)
		fmt.Print("        set_local $tmp\n")
	}

	fmt.Print("        get_local $tmp\n")
	fmt.Print("    )\n")

	fmt.Print("    (func $main (export \"_start\")\n")
	fmt.Print("        (call $proc_exit\n")
	fmt.Print("            (call $__)\n")
	fmt.Print("        )\n")
	fmt.Print("    )\n")
	fmt.Print(")\n")
}
