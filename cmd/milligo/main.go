package main

import (
	"errors"
	"fmt"
	"github.com/sanemat/go-milligo/token"
	"os"
	"strconv"
	"unicode"
)

var errOpMismatch = errors.New("op mismatch")
var errTokenIsNotNum = errors.New("expected a number")
var errTokenizeInt = errors.New("expect number string")
var errInvalidToken = errors.New("invalid token")
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
		return errOpMismatch
	}
	tk = tk.Next
	return nil
}

// Ensure that the current tk is NUM.
func expectNumber() (int, error) {
	if tk.Kind != token.NUM {
		return 0, errTokenIsNotNum
	}
	val := tk.Val
	tk = tk.Next
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
	return tk.Kind == token.EOF
}

func tokenize(s string) (*token.Token, error) {
	head := token.Token{}
	cur := &head
	for i:= 0; i < len(s); i++ {

		// Punctuator
		if string(s[i]) == "+" || string(s[i]) == "-" {
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
				return nil, errTokenizeInt
			}
			cur.Val = n
			i = j - 1
			continue
		}
		return nil, errInvalidToken
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
	if err2 != nil {
		fmt.Print(err2.Error())
		os.Exit(1)
	}

	fmt.Print("(module\n")
	fmt.Print("    (import \"wasi_unstable\" \"proc_exit\" (func $proc_exit (param i32)))\n")
	fmt.Print("    (memory 1)\n")
	fmt.Print("    (export \"memory\" (memory 0))\n")
	fmt.Print("    (func $main (export \"_start\")\n")
	fmt.Print("        (call $proc_exit\n")

	// The first token must be a number
	n, err := expectNumber()
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	fmt.Printf("            i32.const %d\n", n)

	for !atEOF() {
		if consume("+") {
			n, err := expectNumber()
			if err != nil {
				fmt.Print(err.Error())
				os.Exit(1)
			}

			fmt.Printf("            i32.const %d\n", n)
			fmt.Print("            i32.add\n")
			continue
		}

		if err := expect("-"); err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
		n, err := expectNumber()
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		fmt.Printf("            i32.const %d\n", n)
		fmt.Print("            i32.sub\n")
		continue
	}
	fmt.Print("        )\n")
	fmt.Print("    )\n")
	fmt.Print(")\n")
}
