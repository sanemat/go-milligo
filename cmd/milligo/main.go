package main

import (
	"fmt"
	"github.com/sanemat/go-milligo/token"
	"os"
	"strconv"
	"unicode"
)

var tk *token.Token
var userInput string

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
		return fmt.Errorf("%s\nexpected=%s, actual=%s", userInput, op, string(tk.Str[0]))
	}
	tk = tk.Next
	return nil
}

// Ensure that the current tk is NUM.
func expectNumber() (int, error) {
	if tk.Kind != token.NUM {
		return 0, fmt.Errorf("%s\nexpect NUM, got=%s", userInput, tk.Kind)
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

func tokenize() (*token.Token, error) {
	s := userInput
	head := token.Token{}
	cur := &head
	for i:= 0; i < len(s); i++ {
		// Skip whitespace characters.
		if unicode.IsSpace(rune(s[i])) {
			continue
		}

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
				return nil, fmt.Errorf("%s\n%*s^ %s", userInput, i, "", "expect number string")
			}
			cur.Val = n
			i = j - 1
			continue
		}
		return nil, fmt.Errorf("%s\n%*s^ %s", userInput, i, "", "invalid token")
	}
	newToken(token.EOF, cur, "")
	return head.Next, nil
}

var err2 error

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "args is not 2. got=%d", len(os.Args))
		os.Exit(1)
	}

	userInput = os.Args[1]
	tk, err2 = tokenize()
	if err2 != nil {
		fmt.Fprint(os.Stderr, err2.Error())
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
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Printf("            i32.const %d\n", n)

	for !atEOF() {
		if consume("+") {
			n, err := expectNumber()
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
				os.Exit(1)
			}

			fmt.Printf("            i32.const %d\n", n)
			fmt.Print("            i32.add\n")
			continue
		}

		if err := expect("-"); err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
		n, err := expectNumber()
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
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
