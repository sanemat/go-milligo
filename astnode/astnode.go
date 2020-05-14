package astnode

// Kind for Astnode
//go:generate stringer -type=Kind
type Kind int

// Astnode for AST node type
type Astnode struct {
	Kind Kind     // Astnode Kind
	LHS  *Astnode // Left-hand side
	RHS  *Astnode // Right-hand side
	Val  int      // Used if kind == NUM
}

// List of kind
const (
	ADD Kind = iota // +
	SUB             // -
	MUL             // *
	DIV             // /
	EQ              // ==
	NE              // !=
	LT              // <
	LE              // <=
	NUM             // Integer
)
