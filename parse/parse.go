package parse

import (
	"fmt"
	"github.com/sanemat/go-milligo/astnode"
	"github.com/sanemat/go-milligo/tokenize"
	"os"
)

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

// Expr expr = equality
func Expr() *astnode.Astnode {
	return equality()
}

// equality   = relational ("==" relational | "!=" relational)*
func equality() *astnode.Astnode {
	node := relational()
	for {
		if tokenize.Consume("==") {
			node = newBinary(astnode.EQ, node, relational())
		} else if tokenize.Consume("!=") {
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
		if tokenize.Consume("<") {
			node = newBinary(astnode.LT, node, add())
		} else if tokenize.Consume("<=") {
			node = newBinary(astnode.LE, node, add())
		} else if tokenize.Consume(">") {
			node = newBinary(astnode.LT, add(), node)
		} else if tokenize.Consume(">=") {
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
		if tokenize.Consume("+") {
			node = newBinary(astnode.ADD, node, mul())
		} else if tokenize.Consume("-") {
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
		if tokenize.Consume("*") {
			node = newBinary(astnode.MUL, node, unary())
		} else if tokenize.Consume("/") {
			node = newBinary(astnode.DIV, node, unary())
		} else {
			return node
		}
	}
}

// unary = ("+" | "-")? unary
//       | primary
func unary() *astnode.Astnode {
	if tokenize.Consume("+") {
		return unary()
	}
	if tokenize.Consume("-") {
		return newBinary(astnode.SUB, newNum(0), unary())
	}
	return primary()
}

// primary = "(" expr ")" | num
func primary() *astnode.Astnode {
	if tokenize.Consume("(") {
		node := Expr()
		tokenize.Expect(")")
		return node
	}

	n, err := tokenize.ExpectNumber()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
	return newNum(n)
}
