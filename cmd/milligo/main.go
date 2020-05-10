package main

import (
	"fmt"
	"github.com/sanemat/go-milligo"
	"github.com/sanemat/go-milligo/parser"
	"github.com/sanemat/go-milligo/token"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// Consumes the current milligo.Tk if it matches `op`.
func consume(op string) bool {
	if milligo.Tk.Kind != token.RESERVED || milligo.Tk.Str != op {
		return false
	}
	milligo.Tk = milligo.Tk.Next
	return true
}

// Ensure that the current milligo.Tk is `op`.
func expect(op string) error {
	if milligo.Tk.Kind != token.RESERVED || milligo.Tk.Str != op {
		return fmt.Errorf("%s\nexpected=%s, actual=%s", milligo.UserInput, op, milligo.Tk.Str)
	}
	milligo.Tk = milligo.Tk.Next
	return nil
}

// Ensure that the current milligo.Tk is NUM.
func expectNumber() (int, error) {
	if milligo.Tk.Kind != token.NUM {
		return 0, fmt.Errorf("%s\nexpect NUM, got=%s", milligo.UserInput, milligo.Tk.Kind)
	}
	val := milligo.Tk.Val
	milligo.Tk = milligo.Tk.Next
	return val, nil
}

func newToken(kind token.Kind, cur *token.Token, str string) *token.Token {
	tok := token.Token{
		Kind: kind,
		Str:  str,
	}
	cur.Next = &tok
	return &tok
}

func atEOF() bool {
	return milligo.Tk.Kind == token.EOF
}

func tokenize() (*token.Token, error) {
	s := milligo.UserInput
	head := token.Token{}
	cur := &head
	for i := 0; i < len(s); i++ {
		// Skip whitespace characters.
		if unicode.IsSpace(rune(s[i])) {
			continue
		}

		// Multi-letter punctuator
		if len(s)-(i+1) > 0 && (s[i:i+2] == "==" || s[i:i+2] == "!=" || s[i:i+2] == "<=" || s[i:i+2] == ">=") {
			cur = newToken(token.RESERVED, cur, s[i:i+2])
			i++
			continue
		}

		// Single-letter punctuator
		if strings.ContainsAny(string(s[i]), "-+*/()<>") {
			cur = newToken(token.RESERVED, cur, string(s[i]))
			continue
		}

		// Integer literal
		if unicode.IsDigit(rune(s[i])) {
			var j int
			for j = i + 1; j < len(s); j++ {
				if !unicode.IsDigit(rune(s[j])) {
					break
				}
			}
			cur = newToken(token.NUM, cur, s[i:j])
			n, err := strconv.Atoi(s[i:j])
			if err != nil {
				return nil, fmt.Errorf("%s\n%*s^ %s", milligo.UserInput, i, "", "expect number string")
			}
			cur.Val = n
			i = j - 1
			continue
		}
		return nil, fmt.Errorf("%s\n%*s^ %s", milligo.UserInput, i, "", "invalid token")
	}
	newToken(token.EOF, cur, "")
	return head.Next, nil
}

//
// Parser
//

func newNode(kind parser.NodeKind) parser.Node {
	return parser.Node{
		Kind: kind,
	}
}

func newBinary(kind parser.NodeKind, lhs *parser.Node, rhs *parser.Node) *parser.Node {
	node := newNode(kind)
	node.LHS = lhs
	node.RHS = rhs
	return &node
}

func newNum(val int) *parser.Node {
	node := newNode(parser.NUM)
	node.Val = val
	return &node
}

// expr = equality
func expr() *parser.Node {
	return equality()
}

// equality   = relational ("==" relational | "!=" relational)*
func equality() *parser.Node {
	node := relational()
	for {
		if consume("==") {
			node = newBinary(parser.EQ, node, relational())
		} else if consume("!=") {
			node = newBinary(parser.NE, node, relational())
		} else {
			return node
		}
	}
}

// relational = add ("<" add | "<=" add | ">" add | ">=" add)*
func relational() *parser.Node {
	node := add()

	for {
		if consume("<") {
			node = newBinary(parser.LT, node, add())
		} else if consume("<=") {
			node = newBinary(parser.LE, node, add())
		} else if consume(">") {
			node = newBinary(parser.LT, add(), node)
		} else if consume(">=") {
			node = newBinary(parser.LE, add(), node)
		} else {
			return node
		}
	}
}

// add = mul ("+" mul | "-" mul)*
func add() *parser.Node {
	node := mul()
	for {
		if consume("+") {
			node = newBinary(parser.ADD, node, mul())
		} else if consume("-") {
			node = newBinary(parser.SUB, node, mul())
		} else {
			return node
		}
	}
}

// mul = unary ("*" unary | "/" unary)*
func mul() *parser.Node {
	node := unary()
	for {
		if consume("*") {
			node = newBinary(parser.MUL, node, unary())
		} else if consume("/") {
			node = newBinary(parser.DIV, node, unary())
		} else {
			return node
		}
	}
}

// unary = ("+" | "-")? unary
//       | primary
func unary() *parser.Node {
	if consume("+") {
		return unary()
	}
	if consume("-") {
		return newBinary(parser.SUB, newNum(0), unary())
	}
	return primary()
}

// primary = "(" expr ")" | num
func primary() *parser.Node {
	if consume("(") {
		node := expr()
		expect(")")
		return node
	}

	n, err := expectNumber()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
	return newNum(n)
}

//
// Code generator
//
func gen(node *parser.Node) {
	if node.Kind == parser.NUM {
		fmt.Printf("            i32.const %d\n", node.Val)
		return
	}

	gen(node.LHS)
	gen(node.RHS)

	switch node.Kind {
	case parser.ADD:
		fmt.Print("            i32.add\n")
	case parser.SUB:
		fmt.Print("            i32.sub\n")
	case parser.MUL:
		fmt.Print("            i32.mul\n")
	case parser.DIV:
		fmt.Print("            i32.div_s\n")
	case parser.EQ:
		fmt.Print("            i32.eq\n")
	case parser.NE:
		fmt.Print("            i32.ne\n")
	case parser.LT:
		fmt.Print("            i32.lt_s\n")
	case parser.LE:
		fmt.Print("            i32.le_s\n")
	}
}

var err2 error

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "args is not 2. got=%d", len(os.Args))
		os.Exit(1)
	}

	// Tokenize and parse
	milligo.UserInput = os.Args[1]
	milligo.Tk, err2 = tokenize()
	if err2 != nil {
		fmt.Fprint(os.Stderr, err2.Error())
		os.Exit(1)
	}
	node := expr()

	fmt.Print("(module\n")
	fmt.Print("    (import \"wasi_unstable\" \"proc_exit\" (func $proc_exit (param i32)))\n")
	fmt.Print("    (memory 1)\n")
	fmt.Print("    (export \"memory\" (memory 0))\n")
	fmt.Print("    (func $main (export \"_start\")\n")
	fmt.Print("        (call $proc_exit\n")

	// Traverse the AST to emit assembly.
	gen(node)

	fmt.Print("        )\n")
	fmt.Print("    )\n")
	fmt.Print(")\n")
}
