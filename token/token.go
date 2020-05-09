package token

// Kind for token
type Kind string

// Token of program
type Token struct {
	Kind Kind
	Next *Token
	Val  int
	Str  string
}

// List of kind
const (
	RESERVED = "RESERVED"
	NUM      = "NUM"
	EOF      = "EOF"
)
