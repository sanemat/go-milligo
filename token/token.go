package token

// Kind for token
//go:generate stringer -type=Kind
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
	RESERVED Kind = iota
	NUM
	EOF
	RETURN
)
