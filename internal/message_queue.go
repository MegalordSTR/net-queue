package internal

import (
	"context"
	"sync"
)

type MessageQueueMap struct {
	q sync.Map
}

func NewMessageQueueMap() *MessageQueueMap {
	return &MessageQueueMap{}
}

func (m *MessageQueueMap) GetQueue(name string) *MessageQueue {
	q, _ := m.q.LoadOrStore(name, NewMessageQueue())
	return q.(*MessageQueue)
}

type MessageQueue struct {
	items chan []string
	empty chan bool
}

func NewMessageQueue() *MessageQueue {
	q := &MessageQueue{
		items: make(chan []string, 1),
		empty: make(chan bool, 1),
	}

	q.empty <- true

	return q
}

func (q *MessageQueue) Put(item string) {
	var items []string
	select {
	case items = <-q.items:
	case <-q.empty:
	}
	items = append(items, item)
	q.items <- items
}

func (q *MessageQueue) GetOrReturn() string {
	var items []string
	select {
	case items = <-q.items:
	default:
		return ""
	}

	return q.get(items)
}

// Очередность взятия(в большинстве случаев) обеспечивается через каналы (согласно описанию мьютекса (канал содержит свой))
func (q *MessageQueue) GetCtx(ctx context.Context) string {
	var items []string
	select {
	case items = <-q.items:
	case <-ctx.Done():
		return ""
	}

	return q.get(items)
}

func (q *MessageQueue) get(items []string) string {
	item := items[0]
	items = items[1:]
	if len(items) == 0 {
		q.empty <- true
	} else {
		q.items <- items
	}
	return item
}
