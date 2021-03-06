package git

/*
Copied from git c code
*/
import (
	"fmt"
	"strings"

	"github.com/liy/goe/object"
)

type CommitRank struct {
	commit       *object.Commit
	enqueueIndex int
}

func (cr *CommitRank) GetRank() int {
	if cr.commit.Committer != cr.commit.Author {
		return int(cr.commit.Committer.TimeStamp.Unix())
	}
	return int(cr.commit.Author.TimeStamp.Unix())
}

type CommitPrioQueue struct {
	queue []CommitRank
}

func (q CommitPrioQueue) String() string {
	var sb strings.Builder
	fmt.Fprint(&sb, "[")
	for _, item := range q.queue {
		fmt.Fprintf(&sb, "%v ", item.GetRank())
	}
	fmt.Fprint(&sb, "]")
	return sb.String()
}

func (q *CommitPrioQueue) Enqueue(commit *object.Commit) {
	q.queue = append(q.queue, CommitRank{
		commit,
		len(q.queue),
	})

	var child int
	for i := len(q.queue) - 1; i > 0; i = child {
		child = (i - 1) / 2
		if q.compare(child, i) < 0 {
			break
		}

		// swap
		q.queue[child], q.queue[i] = q.queue[i], q.queue[child]
	}
}

func (q *CommitPrioQueue) Dequeue() *object.Commit {
	if len(q.queue) == 0 {
		return nil
	}

	result := q.queue[0]
	if len(q.queue) == 0 {
		return result.commit
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

		// swap
		q.queue[child], q.queue[i] = q.queue[i], q.queue[child]
	}

	return result.commit
}

func (q *CommitPrioQueue) compare(a int, b int) int {
	if q.queue[a].GetRank() < q.queue[b].GetRank() {
		return 1
	} else if q.queue[a].GetRank() > q.queue[b].GetRank() {
		return -1
	} else {
		return q.queue[a].enqueueIndex - q.queue[b].enqueueIndex
	}
}

func (q *CommitPrioQueue) Size() int {
	return len(q.queue)
}

func (q *CommitPrioQueue) ForEach(callback func(CommitRank, int)) {
	for i, o := range q.queue {
		callback(o, i)
	}
}