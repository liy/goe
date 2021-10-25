package utils

/*
Copied from git c code
*/
import (
	"fmt"
	"strings"
)

type Comparable interface {
	GetCompareValue() int
}

type PrioQueueItem struct {
	object       *Comparable
	enqueueIndex int
}

type PrioQueue struct {
	queue []*PrioQueueItem
}

func (q PrioQueue) String() string {
	var sb strings.Builder
	fmt.Fprint(&sb, "[")
	for _, item := range q.queue {
		fmt.Fprintf(&sb, "%v ", (*(item.object)).GetCompareValue())
	}
	fmt.Fprint(&sb, "]")
	return sb.String()
}

func (q *PrioQueue) Enqueue(object Comparable) {
	item := &PrioQueueItem{
		&object,
		len(q.queue),
	}

	q.queue = append(q.queue, item)

	var child int
	for i := len(q.queue) - 1; i > 0; i = child {
		child = (i - 1) / 2
		if q.compare(child, i) < 0 {
			break
		}

		q.swap(child, i)
	}
}

func (q *PrioQueue) Dequeue() *Comparable {
	if len(q.queue) == 0 {
		return nil
	}

	result := q.queue[0]
	if len(q.queue) == 0 {
		return result.object
	}

	q.queue[0] = q.queue[len(q.queue)-1]
	// Remove first element
	q.queue = q.queue[:len(q.queue)-1]

	var child int
	l := len(q.queue)
	for i := 0; i*2+1 < l; i = child {
		child = i*2 + 1
		if child+1 < l && q.compare(child, child+1) >= 0 {
			child++
		}

		if q.compare(i, child) <= 0 {
			break
		}

		q.swap(child, i)
	}

	return result.object
}

func (q *PrioQueue) compare(a int, b int) int {
	if (*(q.queue[a].object)).GetCompareValue() < (*(q.queue[b].object)).GetCompareValue() {
		return 1
	} else if (*(q.queue[a].object)).GetCompareValue() > (*(q.queue[b].object)).GetCompareValue() {
		return -1
	} else {
		return q.queue[a].enqueueIndex - q.queue[b].enqueueIndex
	}
}

func (q *PrioQueue) swap(a int, b int) {
	t := q.queue[a]
	q.queue[a] = q.queue[b]
	q.queue[b] = t
}

func (q *PrioQueue) Size() int {
	return len(q.queue)
}

func (q *PrioQueue) ForEach(callback func(*Comparable, int)) {
	for i, o := range q.queue {
		callback(o.object, i)
	}
}