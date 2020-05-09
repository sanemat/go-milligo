package token

type Kind string

type Token struct {
	Kind Kind
	Next *Token
	Val  int
	Str  string
}

const (
	RESERVED = "RESERVED"
	NUM      = "NUM"
	EOF      = "EOF"
)
