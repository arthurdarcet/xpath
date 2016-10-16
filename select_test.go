package gxpath

import (
	"bytes"
	"strings"
	"testing"

	"github.com/antchfx/gxpath/xpath"
)

var html *TNode = example()

/*
testXPath2(t, html, "//title/node()", 1) still have some error,will fix in future.
testXPath(t, html, "//*[count(*)=3]", "ul")
testXPath3(t, html, "//li[floor(3 div 2)]", selectNode(html, "//li[1]"))
*/

func TestSelf(t *testing.T) {
	testXPath(t, html, ".", "html")
	testXPath(t, html.FirstChild, ".", "head")
	testXPath(t, html, "self::*", "html")
	testXPath(t, html.LastChild, "self::body", "body")
	testXPath2(t, html, "//body/./ul/li/a", 3)
}

func TestParent(t *testing.T) {
	testXPath(t, html.LastChild, "..", "html")
	testXPath(t, html.LastChild, "parent::*", "html")
	a := selectNode(html, "//li/a")
	testXPath(t, a, "parent::*", "li")
	testXPath(t, html, "//title/parent::head", "head")
}

func TestAttribute(t *testing.T) {
	testXPath(t, html, "@lang='en'", "html")
	testXPath2(t, html, "@lang='zh'", 0)
	testXPath2(t, html, "//@href", 3)
	testXPath2(t, html, "//a[@*]", 3)
}

func TestRelativePath(t *testing.T) {
	testXPath(t, html, "head", "head")
	testXPath(t, html, "/head", "head")

	testXPath(t, html, "/head/title", "title")

	testXPath2(t, html, "/body/ul/li/a", 3)
	testXPath(t, html, "//title", "title")
	testXPath(t, html, "//title/..", "head")
	testXPath(t, html, "//title/../..", "html")
	testXPath2(t, html, "//a[@href]", 3)
	testXPath(t, html, "//ul/../footer", "footer")
}

func TestChild(t *testing.T) {
	testXPath(t, html, "/child::head", "head")
	testXPath(t, html, "/child::head/child::title", "title")
	testXPath(t, html, "//title/../child::title", "title")
	testXPath(t, html.Parent, "//child::*", "html")
}

func TestDescendant(t *testing.T) {
	testXPath2(t, html, "descendant::*", 15)
	testXPath2(t, html, "/head/descendant::*", 2)
	testXPath2(t, html, "//ul/descendant::*", 7)  // <li> + <a>
	testXPath2(t, html, "//ul/descendant::li", 4) // <li>
}

func TestAncestor(t *testing.T) {
	testXPath2(t, html, "/body/footer/ancestor::*", 2) // body>html
	testXPath2(t, html, "/body/ul/li/a/ancestor::li", 3)
}

func TestFollowingSibling(t *testing.T) {
	var list []*TNode
	list = selectNodes(html, "//li/following-sibling::*")
	for _, n := range list {
		if n.Data != "li" {
			t.Fatalf("expected node is li,but got:%s", n.Data)
		}
	}

	list = selectNodes(html, "//ul/following-sibling::*") // p,footer
	for _, n := range list {
		if n.Data != "p" && n.Data != "footer" {
			t.Fatal("expected node is not one of the following nodes: [p,footer]")
		}
	}
	testXPath(t, html, "//ul/following-sibling::footer", "footer")
}

func TestPrecedingSibling(t *testing.T) {
	testXPath(t, html, "/body/footer/preceding-sibling::*", "p")
	testXPath2(t, html, "/body/footer/preceding-sibling::*", 3) // p,ul,h1
}

func TestFollowing(t *testing.T) {
	//testXPath2(t, html, "//ul/following:*", 9)
}

func TestStar(t *testing.T) {
	testXPath(t, html, "/head/*", "title")
	testXPath2(t, html, "//ul/*", 4)
	testXPath(t, html, "@*", "html")
	testXPath2(t, html, "/body/h1/*", 0)
	testXPath2(t, html, `//ul/*/a`, 3)
}

func TestNodeTestType(t *testing.T) {
	testXPath(t, html, "//title/text()", "Hello")
	testXPath(t, html, "//a[@href='/']/text()", "Home")
	testXPath2(t, html, "//head/node()", 2)

}

func TestPosition(t *testing.T) {
	testXPath3(t, html, "/head[1]", html.FirstChild) // compare to 'head' element
	ul := selectNode(html, "//ul")
	testXPath3(t, html, "/head[last()]", html.FirstChild)
	testXPath3(t, html, "//li[1]", ul.FirstChild)
	testXPath3(t, html, "//li[4]", ul.LastChild)
	testXPath3(t, html, "//li[last()]", ul.LastChild)
}

func TestPredicate(t *testing.T) {
	testXPath(t, html.Parent, "html[@lang='en']", "html")
	testXPath(t, html, "//a[@href='/']", "a")
	testXPath(t, html, "//meta[@name]", "meta")
	ul := selectNode(html, "//ul")
	testXPath3(t, html, "//li[position()=4]", ul.LastChild)
	testXPath3(t, html, "//li[position()=1]", ul.FirstChild)
	testXPath2(t, html, "//li[position()>0]", 4)
}

func TestFunction(t *testing.T) {
	// count(Node-Set)
	//
}

func TestOperationOrLogical(t *testing.T) {
	testXPath3(t, html, "//li[1+1]", selectNode(html, "//li[2]"))
	testXPath3(t, html, "//li[4 div 2]", selectNode(html, "//li[2]"))
	testXPath3(t, html, "//li[3 mod 2]", selectNode(html, "//li[1]"))
	testXPath3(t, html, "//li[3 - 2]", selectNode(html, "//li[1]"))
	testXPath2(t, html, "//a[@id>=1]", 3)         // //a[@id>=1] == a[1],a[2],a[3]
	testXPath2(t, html, "//a[@id<2]", 1)          // //a[@id>=1] == a[1]
	testXPath2(t, html, "//a[@id!=2]", 2)         // //a[@id>=1] == a[1],a[3]
	testXPath2(t, html, "//a[@id=1 or @id=3]", 2) // //a[@id>=1] == a[1],a[3]
	testXPath3(t, html, "//a[@id=1 and @href='/']", selectNode(html, "//a[1]"))
}

func testXPath(t *testing.T, root *TNode, expr string, expected string) {
	node := selectNode(root, expr)
	if node == nil {
		t.Fatalf("`%s` returns node is nil", expr)
	}
	if node.Data != expected {
		t.Fatalf("`%s` expected node is %s,but got %s", expr, expected, node.Data)
	}
}

func testXPath2(t *testing.T, root *TNode, expr string, expected int) {
	list := selectNodes(root, expr)
	if len(list) != expected {
		t.Fatalf("`%s` expected node numbers is %d,but got %d", expr, expected, len(list))
	}
}

func testXPath3(t *testing.T, root *TNode, expr string, expected *TNode) {
	node := selectNode(root, expr)
	if node == nil {
		t.Fatalf("`%s` returns node is nil", expr)
	}
	if node != expected {
		t.Fatalf("`%s` %s != %s", expr, node.Value(), expected.Value())
	}
}

func selectNode(root *TNode, expr string) (n *TNode) {
	t := Select(createNavigator(root), expr)
	if t.MoveNext() {
		n = (t.Current().(*TNodeNavigator)).curr
	}
	return n
}

func selectNodes(root *TNode, expr string) []*TNode {
	t := Select(createNavigator(root), expr)
	var nodes []*TNode
	for t.MoveNext() {
		node := (t.Current().(*TNodeNavigator)).curr
		nodes = append(nodes, node)
	}
	return nodes
}

func createNavigator(n *TNode) *TNodeNavigator {
	return &TNodeNavigator{curr: n, root: n, attr: -1}
}

type Attribute struct {
	Key, Value string
}

type TNode struct {
	Parent, FirstChild, LastChild, PrevSibling, NextSibling *TNode

	Type xpath.NodeType
	Data string
	Attr []Attribute
}

func (n *TNode) Value() string {
	if n.Type == xpath.TextNode {
		return n.Data
	}

	var buff bytes.Buffer
	var output func(*TNode)
	output = func(node *TNode) {
		if node.Type == xpath.TextNode {
			buff.WriteString(node.Data)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			output(child)
		}
	}
	output(n)
	return buff.String()
}

// TNodeNavigator is for navigating TNode.
type TNodeNavigator struct {
	curr, root *TNode
	attr       int
}

func (n *TNodeNavigator) NodeType() xpath.NodeType {
	if n.curr.Type == xpath.ElementNode && n.attr != -1 {
		return xpath.AttributeNode
	}
	return n.curr.Type
}

func (n *TNodeNavigator) LocalName() string {
	if n.attr != -1 {
		return n.curr.Attr[n.attr].Key
	}
	return n.curr.Data
}

func (n *TNodeNavigator) Prefix() string {
	return ""
}

func (n *TNodeNavigator) Value() string {
	switch n.curr.Type {
	case xpath.CommentNode:
		return n.curr.Data
	case xpath.ElementNode:
		if n.attr != -1 {
			return n.curr.Attr[n.attr].Value
		}
		var buf bytes.Buffer
		node := n.curr.FirstChild
		for node != nil {
			if node.Type == xpath.TextNode {
				buf.WriteString(strings.TrimSpace(node.Data))
			}
			node = node.NextSibling
		}
		return buf.String()
	case xpath.TextNode:
		return n.curr.Data
	}
	return ""
}

func (n *TNodeNavigator) Copy() xpath.NodeNavigator {
	n2 := *n
	return &n2
}

func (n *TNodeNavigator) Current() xpath.Node {
	return n.curr
}

func (n *TNodeNavigator) MoveToRoot() {
	n.curr = n.root
}

func (n *TNodeNavigator) MoveToParent() bool {
	if node := n.curr.Parent; node != nil {
		n.curr = node
		return true
	}
	return false
}

func (n *TNodeNavigator) MoveToNextAttribute() bool {
	if n.attr >= len(n.curr.Attr)-1 {
		return false
	}
	n.attr++
	return true
}

func (n *TNodeNavigator) MoveToChild() bool {
	if node := n.curr.FirstChild; node != nil {
		n.curr = node
		return true
	}
	return false
}

func (n *TNodeNavigator) MoveToFirst() bool {
	if n.curr.PrevSibling == nil {
		return false
	}
	for {
		node := n.curr.PrevSibling
		if node == nil {
			break
		}
		n.curr = node
	}
	return true
}

func (n *TNodeNavigator) String() string {
	return n.Value()
}

func (n *TNodeNavigator) MoveToNext() bool {
	if node := n.curr.NextSibling; node != nil {
		n.curr = node
		return true
	}
	return false
}

func (n *TNodeNavigator) MoveToPrevious() bool {
	if node := n.curr.PrevSibling; node != nil {
		n.curr = node
		return true
	}
	return false
}

func (n *TNodeNavigator) MoveTo(other xpath.NodeNavigator) bool {
	node, ok := other.(*TNodeNavigator)
	if !ok || node.root != n.root {
		return false
	}

	n.curr = node.curr
	n.attr = node.attr
	return true
}

func createNode(data string, typ xpath.NodeType) *TNode {
	return &TNode{Data: data, Type: typ, Attr: make([]Attribute, 0)}
}

func (t *TNode) createChildNode(data string, typ xpath.NodeType) *TNode {
	n := createNode(data, typ)
	n.Parent = t
	if t.FirstChild == nil {
		t.FirstChild = n
	} else {
		t.LastChild.NextSibling = n
		n.PrevSibling = t.LastChild
	}
	t.LastChild = n
	return n
}

func (t *TNode) appendNode(data string, typ xpath.NodeType) *TNode {
	n := createNode(data, typ)
	n.Parent = t.Parent
	t.NextSibling = n
	n.PrevSibling = t
	if t.Parent != nil {
		t.Parent.LastChild = n
	}
	return n
}

func (t *TNode) addAttribute(k, v string) {
	t.Attr = append(t.Attr, Attribute{k, v})
}

func example() *TNode {
	/*
		<html lang="en">
		   <head>
			   <title>Hello</title>
			   <meta name="language" content="en"/>
		   </head>
		   <body>
				<h1>This is a H1</h1>
				<ul>
					<li><a id="1" href="/">Home</a></li>
					<li><a id="2" href="/about">about</a></li>
					<li><a id="3" href="/account">login</a></li>
					<li></li>
				</ul>
				<p>
					Hello,This is an example for gxpath.
				</p>
				<footer>footer script</footer>
		   </body>
		</html>
	*/
	doc := createNode("", xpath.RootNode)
	xhtml := doc.createChildNode("html", xpath.ElementNode)
	xhtml.addAttribute("lang", "en")

	// The HTML head section.
	head := xhtml.createChildNode("head", xpath.ElementNode)
	n := head.createChildNode("title", xpath.ElementNode)
	n = n.createChildNode("Hello", xpath.TextNode)
	n = head.createChildNode("meta", xpath.ElementNode)
	n.addAttribute("name", "language")
	n.addAttribute("content", "en")
	// The HTML body section.
	body := xhtml.createChildNode("body", xpath.ElementNode)
	n = body.createChildNode("h1", xpath.ElementNode)
	n = n.createChildNode("Hello", xpath.TextNode)
	ul := body.createChildNode("ul", xpath.ElementNode)
	n = ul.createChildNode("li", xpath.ElementNode)
	n = n.createChildNode("a", xpath.ElementNode)
	n.addAttribute("id", "1")
	n.addAttribute("href", "/")
	n = n.createChildNode("Home", xpath.TextNode)
	n = ul.createChildNode("li", xpath.ElementNode)
	n = n.createChildNode("a", xpath.ElementNode)
	n.addAttribute("id", "2")
	n.addAttribute("href", "/about")
	n = n.createChildNode("about", xpath.TextNode)
	n = ul.createChildNode("li", xpath.ElementNode)
	n = n.createChildNode("a", xpath.ElementNode)
	n.addAttribute("id", "3")
	n.addAttribute("href", "/account")
	n = n.createChildNode("login", xpath.TextNode)
	n = ul.createChildNode("li", xpath.ElementNode)

	n = body.createChildNode("p", xpath.ElementNode)
	n = n.createChildNode("Hello,This is an example for gxpath.", xpath.TextNode)

	n = body.createChildNode("footer", xpath.ElementNode)
	n = n.createChildNode("footer script", xpath.TextNode)

	return xhtml
}
