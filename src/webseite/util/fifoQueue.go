package util

import "webseite/models"

// Queue is a basic FIFO queue based on a circular list that resizes as needed.
type Queue struct {
    Nodes    []*models.Ping
    head    int
    tail    int
    count    int
}

// Push adds a node to the queue.
func (q *Queue) Push(n *models.Ping) {
    if q.head == q.tail && q.count > 0 {
        nodes := make([]*models.Ping, len(q.Nodes)*2)
        copy(nodes, q.Nodes[q.head:])
        copy(nodes[len(q.Nodes)-q.head:], q.Nodes[:q.head])
        q.head = 0
        q.tail = len(q.Nodes)
        q.Nodes = nodes
    }

    q.Nodes[q.tail] = n
    q.tail = (q.tail + 1) % len(q.Nodes)
    q.count++
}

func (q *Queue) Size() int {
    return q.count
}

// Pop removes and returns a node from the queue in first to last order.
func (q *Queue) Pop() *models.Ping {
    if q.count == 0 {
        return nil
    }

    node := q.Nodes[q.head]
    q.head = (q.head + 1) % len(q.Nodes)
    q.count--
    return node
}