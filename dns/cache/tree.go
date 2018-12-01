package cache

import (
	"errors"
)

// Node is a struct that is used to store a node (including it's contents) and a map[string] of pointers to it's children.
type Node struct {
	children map[string]*Node
	Content  interface{}
}

// GetChild gets the pointer to a specific child of a node
func (n Node) GetChild(name string) *Node {
	c, ok := n.children[name]
	if !ok {
		return nil
	}
	return c
}

// AddChild adds a child to a node
func (n *Node) AddChild(name string, c *Node) error {
	if _, ok := n.children[name]; ok {
		return errors.New("child does already exist")
	}
	n.children[name] = c
	return nil
}

// RemoveChild removes the specified child from a node
func (n *Node) RemoveChild(name string) error {
	if _, ok := n.children[name]; !ok {
		return errors.New("child does not exist")
	}
	delete(n.children, name)
	return nil
}
