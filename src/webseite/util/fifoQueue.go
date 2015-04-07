package util

import (
    "fmt"
)

type Node struct {
    Value int
}

// Stack is a basic LIFO stack that resizes as needed.
type Stack struct {
    nodes    []*Node
    count    int
}

// Push adds a node to the stack.
func (s *Stack) Push(n *Node) {
    if s.count >= len(s.nodes) {
        nodes := make([]*Node, len(s.nodes)*2)
        copy(nodes, s.nodes)
        s.nodes = nodes
    }

    s.nodes[s.count] = n
    s.count++
}

// Pop removes and returns a node from the stack in last to first order.
func (s *Stack) Pop() *Node {
    if s.count == 0 {
        return nil
    }

    node := s.nodes[s.count-1]
    s.count--
    return node
}

// Queue is a basic FIFO queue based on a circular list that resizes as needed.
type Queue struct {
    nodes    []*Node
    head    int
    tail    int
    count    int
}

// Push adds a node to the queue.
func (q *Queue) Push(n *Node) {
    if q.head == q.tail && q.count > 0 {
        nodes := make([]*Node, len(q.nodes)*2)
        copy(nodes, q.nodes[q.head:])
        copy(nodes[len(q.nodes)-q.head:], q.nodes[:q.head])
        q.head = 0
        q.tail = len(q.nodes)
        q.nodes = nodes
    }

    q.nodes[q.tail] = n
    q.tail = (q.tail + 1) % len(q.nodes)
    q.count++
}

func (q *Queue) Size() int {
    return q.count
}

// Pop removes and returns a node from the queue in first to last order.
func (q *Queue) Pop() *Node {
    if q.count == 0 {
        return nil
    }

    node := q.nodes[q.head]
    q.head = (q.head + 1) % len(q.nodes)
    q.count--
    return node
}