package parser

// NodeKind for Node
type NodeKind string

// Node for AST node type
type Node struct {
	Kind NodeKind // Node NodeKind
	LHS  *Node    // Left-hand side
	RHS  *Node    // Right-hand side
	Val  int      // Used if kind == NUM
}

// List of kind
const (
	ADD = "ADD" // +
	SUB = "SUB" // -
	MUL = "MUL" // *
	DIV = "DIV" // /
	NUM = "NUM" // Integer
)
