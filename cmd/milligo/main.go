package main

import (
	"fmt"
	"github.com/sanemat/go-milligo"
	"github.com/sanemat/go-milligo/astnode"
	"github.com/sanemat/go-milligo/codegen"
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

func newNode(kind astnode.Kind) astnode.Astnode {
	return astnode.Astnode{
		Kind: kind,
	}
}

func newBinary(kind astnode.Kind, lhs *astnode.Astnode, rhs *astnode.Astnode) *astnode.Astnode {
	node := newNode(kind)
	node.LHS = lhs
	node.RHS = rhs
	return &node
}

func newNum(val int) *astnode.Astnode {
	node := newNode(astnode.NUM)
	node.Val = val
	return &node
}

// expr = equality
func expr() *astnode.Astnode {
	return equality()
}

// equality   = relational ("==" relational | "!=" relational)*
func equality() *astnode.Astnode {
	node := relational()
	for {
		if consume("==") {
			node = newBinary(astnode.EQ, node, relational())
		} else if consume("!=") {
			node = newBinary(astnode.NE, node, relational())
		} else {
			return node
		}
	}
}

// relational = add ("<" add | "<=" add | ">" add | ">=" add)*
func relational() *astnode.Astnode {
	node := add()

	for {
		if consume("<") {
			node = newBinary(astnode.LT, node, add())
		} else if consume("<=") {
			node = newBinary(astnode.LE, node, add())
		} else if consume(">") {
			node = newBinary(astnode.LT, add(), node)
		} else if consume(">=") {
			node = newBinary(astnode.LE, add(), node)
		} else {
			return node
		}
	}
}

// add = mul ("+" mul | "-" mul)*
func add() *astnode.Astnode {
	node := mul()
	for {
		if consume("+") {
			node = newBinary(astnode.ADD, node, mul())
		} else if consume("-") {
			node = newBinary(astnode.SUB, node, mul())
		} else {
			return node
		}
	}
}

// mul = unary ("*" unary | "/" unary)*
func mul() *astnode.Astnode {
	node := unary()
	for {
		if consume("*") {
			node = newBinary(astnode.MUL, node, unary())
		} else if consume("/") {
			node = newBinary(astnode.DIV, node, unary())
		} else {
			return node
		}
	}
}

// unary = ("+" | "-")? unary
//       | primary
func unary() *astnode.Astnode {
	if consume("+") {
		return unary()
	}
	if consume("-") {
		return newBinary(astnode.SUB, newNum(0), unary())
	}
	return primary()
}

// primary = "(" expr ")" | num
func primary() *astnode.Astnode {
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
	codegen.Codegen(node)
}
