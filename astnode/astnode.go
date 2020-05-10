package astnode

// Kind for AstNode
type Kind string

// AstNode for AST node type
type AstNode struct {
	Kind Kind     // AstNode Kind
	LHS  *AstNode // Left-hand side
	RHS  *AstNode // Right-hand side
	Val  int      // Used if kind == NUM
}

// List of kind
const (
	ADD = "ADD" // +
	SUB = "SUB" // -
	MUL = "MUL" // *
	DIV = "DIV" // /
	EQ  = "EQ"  // ==
	NE  = "NE"  // !=
	LT  = "LT"  // <
	LE  = "LE"  // <=
	NUM = "NUM" // Integer
)
