// ast.go implements an Abstract Syntax Tree for use by the DSL parser.
// The user tells the AST to add nodes and tokens inside the user parse
// function using three basic functions; p.AddNode(), p.AddToken() and
// p.WalkUp(). AST node types are defined by the user.
//
// The AST is made up of nodes, each of which contains a slice of Node
// children and a reference to it's parent.
// The nodes are made available to the user so they can walk up and down
// the tree once it is returned from the parser.
package dsl

import (
	"fmt"
)

// RootNode is the entry point to the tree. curNode is used internally
// to keep track of where the next node should be added.
type AST struct {
	RootNode *Node `json:"root"`
	curNode  *Node `json:"-"`
}

// A Node can contain multiple Tokens which can be useful if the user knows how
// many Tokens belong to a particular Node type. Otherwise, the user should only
// add one token per node.
type Node struct {
	Type     NodeType `json:"type"`
	Tokens   []Token  `json:"tokens"`
	Parent   *Node    `json:"-"`
	Children []Node   `json:"children"`
}

type NodeType string

const (
	NODE_ROOT NodeType = "ROOT"
)

// newAST returns a new instance of AST. The RootNode has the
// builtin node type AST_ROOT.
func newAST() AST {
	rootNode := &Node{Type: NODE_ROOT}
	return AST{RootNode: rootNode, curNode: rootNode}
}

// ---------------------------------------------------------------------------------------------------------

// The AST is made up entirely of Node instances connected in a tree pattern.
// The AST ensures the Parent and Children references are set correctly as it
// is being constructed.
//

// Prints the entire AST tree. It does so by recursively calling Print() on
// each node in the tree in a depth first approach.
//

// Inspect traverses an AST in depth-first order: It starts by calling
// f(node);
func (a *AST) Inspect(fn func(*Node)) {
	visit(a.RootNode, fn)
}

func visit(node *Node, fn func(*Node)) {
	for _, child := range node.Children {
		visit(&child, fn)
	}
	fn(node)
}

// Prints the entire AST tree. It does so by recursively calling Print() on
// each node in the tree in a depth first approach.
func (a *AST) Print() {
	a.RootNode.Print("", true)
	fmt.Println()
}

// Called by Parser.AddNode() in the user parse function. Creates a new node and
// builds the two-way reference to its parent. Also moves the AST curNode
// down the tree to the new node.
func (a *AST) addNode(nt NodeType) {
	a.curNode.Children = append(a.curNode.Children, Node{Type: nt, Parent: a.curNode})
	a.curNode = &a.curNode.Children[len(a.curNode.Children)-1]
}

// Called by Parser.AddToken() in the user parse function. Adds a token to the
// end of the Token slice belonging to the current node.
//
// If Parser.AddToken() is called without any tokens available on the Parser.toks buffer
// the call to AddToken will be logged but no tokens will be added to the node.
func (a *AST) addToken(toks []Token) {
	if toks != nil {
		tokens := append(a.curNode.Tokens, toks...)
		a.curNode.Tokens = tokens
	}
}

// Called by Parser.WalkUp() in the user parse function. Moves the AST
// curNode to its parent.
func (a *AST) walkUp() {
	if a.curNode.Type != NODE_ROOT {
		a.curNode = a.curNode.Parent
	}
}

// Print is a recursive function that keeps track of where in the tree the node
// belongs to so it can print a pretty prefix. The prefix indicates how deep the
// node is and if it is the last node at that level.
//
// A user can print the entire tree using AST.Print() or only print a sub-branch
// by calling Print() on any node in the tree.
func (n *Node) Print(prefix string, isTail bool) {
	fmt.Printf("\n%v", prefix)
	if isTail {
		fmt.Printf("└── ")
	} else {
		fmt.Printf("├── ")
	}
	fmt.Printf("%v - ", n.Type)
	for _, token := range n.Tokens {
		for _, rn := range token.Literal {
			fmt.Print(string(rn))
		}
		fmt.Print(", ")
	}
	numNodes := len(n.Children)
	if numNodes > 0 {
		for _, node := range n.Children[:numNodes-1] {
			if isTail {
				node.Print(prefix+"    ", false)
			} else {
				node.Print(prefix+"│   ", false)
			}
		}
	}
	if numNodes > 0 {
		if isTail {
			n.Children[numNodes-1].Print(prefix+"    ", true)
		} else {
			n.Children[numNodes-1].Print(prefix+"│   ", true)
		}
	}
}
