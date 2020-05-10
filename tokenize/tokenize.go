package tokenize

import (
	"fmt"
	"github.com/sanemat/go-milligo"
	"github.com/sanemat/go-milligo/token"
	"strconv"
	"strings"
	"unicode"
)

// Consume the current milligo.Tk if it matches `op`.
func Consume(op string) bool {
	if milligo.Tk.Kind != token.RESERVED || milligo.Tk.Str != op {
		return false
	}
	milligo.Tk = milligo.Tk.Next
	return true
}

// Expect ensure that the current milligo.Tk is `op`.
func Expect(op string) error {
	if milligo.Tk.Kind != token.RESERVED || milligo.Tk.Str != op {
		return fmt.Errorf("%s\nexpected=%s, actual=%s", milligo.UserInput, op, milligo.Tk.Str)
	}
	milligo.Tk = milligo.Tk.Next
	return nil
}

// ExpectNumber ensure that the current milligo.Tk is NUM.
func ExpectNumber() (int, error) {
	if milligo.Tk.Kind != token.NUM {
		return 0, fmt.Errorf("%s\nexpect NUM, got=%s", milligo.UserInput, milligo.Tk.Kind)
	}
	val := milligo.Tk.Val
	milligo.Tk = milligo.Tk.Next
	return val, nil
}

func atEOF() bool {
	return milligo.Tk.Kind == token.EOF
}

func newToken(kind token.Kind, cur *token.Token, str string) *token.Token {
	tok := token.Token{
		Kind: kind,
		Str:  str,
	}
	cur.Next = &tok
	return &tok
}

// Tokenize milligo.UserInput
func Tokenize() (*token.Token, error) {
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
