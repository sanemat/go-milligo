package token

// Kind for token
type Kind int

// Token of program
type Token struct {
	Kind Kind
	Next *Token
	Val  int
	Str  string
}

// List of kind
const (
	RESERVED = iota
	NUM
	EOF
)
