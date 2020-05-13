package astnode

// Kind for Astnode
type Kind string

// Astnode for AST node type
type Astnode struct {
	Kind Kind     // Astnode Kind
	Next *Astnode // Next node
	LHS  *Astnode // Left-hand side
	RHS  *Astnode // Right-hand side
	Val  int      // Used if kind == NUM
}

// List of kind
const (
	ADD = "ADD"       // +
	SUB = "SUB"       // -
	MUL = "MUL"       // *
	DIV = "DIV"       // /
	EQ  = "EQ"        // ==
	NE  = "NE"        // !=
	LT  = "LT"        // <
	LE  = "LE"        // <=
	NUM = "NUM"       // Integer
	SCO = "SEMICOLON" // ;
)
