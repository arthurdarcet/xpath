package xpath

import (
	"errors"
	"fmt"
)

// NodeType represents a type of XPath node.
type NodeType int

const (
	// RootNode is a root node of the XML document or node tree.
	RootNode NodeType = iota

	// ElementNode is an element, such as <element>.
	ElementNode

	// AttributeNode is an attribute, such as id='123'.
	AttributeNode

	// TextNode is the text content of a node.
	TextNode

	// CommentNode is a comment node, such as <!-- my comment -->
	CommentNode

	// allNode is any types of node, used by xpath package only to predicate match.
	allNode
)

// NodeNavigator provides cursor model for navigating XML data.
type NodeNavigator interface {
	// NodeType returns the XPathNodeType of the current node.
	NodeType() NodeType

	// NamespaceURI gets the namespace uri of the current node.
	NamespaceURI() string

	// LocalName gets the Name of the current node.
	LocalName() string

	// Value gets the value of current node.
	Value() string

	// Copy does a deep copy of the NodeNavigator and all its components.
	Copy() NodeNavigator

	// MoveToRoot moves the NodeNavigator to the root node of the current node.
	MoveToRoot()

	// MoveToParent moves the NodeNavigator to the parent node of the current node.
	MoveToParent() bool

	// MoveToNextAttribute moves the NodeNavigator to the next attribute on current node.
	MoveToNextAttribute() bool

	// MoveToChild moves the NodeNavigator to the first child node of the current node.
	MoveToChild() bool

	// MoveToFirst moves the NodeNavigator to the first sibling node of the current node.
	MoveToFirst() bool

	// MoveToNext moves the NodeNavigator to the next sibling node of the current node.
	MoveToNext() bool

	// MoveToPrevious moves the NodeNavigator to the previous sibling node of the current node.
	MoveToPrevious() bool

	// MoveTo moves the NodeNavigator to the same position as the specified NodeNavigator.
	MoveTo(NodeNavigator) bool
}

// NodeIterator holds all matched Node object.
type NodeIterator struct {
	node      NodeNavigator
	query     query
	ns2prefix func(string) string
}

// Current returns current node which matched.
func (t *NodeIterator) Current() NodeNavigator {
	return t.node
}

// MoveNext moves Navigator to the next match node.
func (t *NodeIterator) MoveNext() bool {
	n := t.query.Select(t)
	if n != nil {
		if !t.node.MoveTo(n) {
			t.node = n.Copy()
		}
		return true
	}
	return false
}

func (t *NodeIterator) NamespaceToPrefix(n string) string {
	return t.ns2prefix(n)
}

// Expr is an XPath expression for query.
type Expr struct {
	s string
	q query
}

// Select selects a node set using the specified XPath expression.
func (expr *Expr) SelectWithNS(root NodeNavigator, ns func(string) string) *NodeIterator {
	return &NodeIterator{query: expr.q.Clone(), node: root, ns2prefix: ns}
}

func (expr *Expr) Select(root NodeNavigator) *NodeIterator {
	return expr.SelectWithNS(root, func(ns string) string { return ns })
}

type EvaluateIterator struct {
	root      NodeNavigator
	ns2prefix func(string) string
}

func (i *EvaluateIterator) Current() NodeNavigator {
	return i.root
}
func (i *EvaluateIterator) NamespaceToPrefix(ns string) string {
	return i.ns2prefix(ns)
}

// Evaluate returns the result of the expression.
// The result type of the expression is one of the follow: bool,float64,string).
func (expr *Expr) EvaluateWithNS(root NodeNavigator, ns func(string) string) interface{} {
	it := EvaluateIterator{root: root, ns2prefix: ns}
	val := expr.q.Evaluate(&it)
	switch val.(type) {
	case query:
		return expr.SelectWithNS(root, ns)
	}
	return val
}

func (expr *Expr) Evaluate(root NodeNavigator) interface{} {
	return expr.EvaluateWithNS(root, func(ns string) string { return ns })
}

// Select selects a node set using the specified XPath expression.
// This method is deprecated, recommend using Expr.Select() method instead.
func Select(root NodeNavigator, expr string) *NodeIterator {
	exp, err := Compile(expr)
	if err != nil {
		panic(err)
	}
	return exp.Select(root)
}

// String returns XPath expression string.
func (expr *Expr) String() string {
	return expr.s
}

// Compile compiles an XPath expression string.
func Compile(expr string) (*Expr, error) {
	if expr == "" {
		return nil, errors.New("expr expression is nil")
	}
	qy, err := build(expr)
	if err != nil {
		return nil, err
	}
	if qy == nil {
		return nil, fmt.Errorf(fmt.Sprintf("undeclared variable in XPath expression: %s", expr))
	}
	return &Expr{s: expr, q: qy}, nil
}

// MustCompile compiles an XPath expression string and ignored error.
func MustCompile(expr string) *Expr {
	exp, err := Compile(expr)
	if err != nil {
		return &Expr{s: expr, q: nopQuery{}}
	}
	return exp
}
