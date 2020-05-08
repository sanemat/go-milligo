package main

import (
	"errors"
	"fmt"
	"github.com/sanemat/go-milligo/token"
	"os"
	"strconv"
	"unicode"
)

var ErrOpMismatch = errors.New("op mismatch")
var ErrTokenIsNotNum = errors.New("expected a number")
var ErrTokenize = errors.New("not integer")
var ErrInvalidToken = errors.New("invalid token")
var tk *token.Token

// Consumes the current tk if it matches `op`.
func consume(op string) bool {
	if tk.Kind != token.RESERVED || string(tk.Str[0]) != op {
		return false
	}
	tk = tk.Next
	return true
}

// Ensure that the current tk is `op`.
func expect(op string) error {
	if tk.Kind != token.RESERVED || string(tk.Str[0]) != op {
		return ErrOpMismatch
	}
	tk = tk.Next
	return nil
}

// Ensure that the current tk is NUM.
func expectNumber() (int, error) {
	if tk.Kind != token.NUM {
		return 0, ErrTokenIsNotNum
	}
	val := tk.Val
	tk = tk.Next
	return val, nil
}

func atEof() bool {
	return tk.Kind == token.EOF
}

func newToken(kind token.Kind, cur *token.Token, str string) *token.Token {
	tok := token.Token{
		Kind: kind,
		Str:  str,
	}
	cur.Next = &tok
	return &tok
}

func tokenize(s string) (*token.Token, error) {
	head := token.Token{}
	cur := &head
	for i:= 0; i <= len(s); i++ {
		// Skip whitespace characters.
		if unicode.IsSpace(rune(s[i])) {
			continue
		}

		// Punctuator
		if rune(s[i]) == '+' || rune(s[i]) == '-' {
			cur = newToken(token.RESERVED, cur, string(s[i]))
			continue
		}

		// Integer literal
		if unicode.IsDigit(rune(s[i])) {
			ns := ""
			j := i
			for {
				if !unicode.IsDigit(rune(s[j])) {
					n, err := strconv.Atoi(ns)
					if err != nil {
						return nil, ErrTokenize
					}
					cur = newToken(token.NUM, cur, ns)
					cur.Val = n
					break
				}
				j++
			}
			continue
		}
		return nil, ErrInvalidToken
	}

	newToken(token.EOF, cur, "")
	return head.Next, nil
}

var err2 error

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("args is not 2. got=%d", len(os.Args))
		os.Exit(1)
	}

	tk, err2 = tokenize(os.Args[1])

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
