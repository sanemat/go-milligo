package astnode

// Kind for AstNode
type Kind int

// AstNode for AST node type
type AstNode struct {
	Kind Kind     // AstNode Kind
	LHS  *AstNode // Left-hand side
	RHS  *AstNode // Right-hand side
	Val  int      // Used if kind == NUM
}

// List of kind
const (
	ADD = iota // +
	SUB        // -
	MUL        // *
	DIV        // /
	EQ         // ==
	NE         // !=
	LT         // <
	LE         // <=
	NUM        // Integer
)
