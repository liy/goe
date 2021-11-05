package git

import (
	"errors"

	"github.com/liy/goe/object"
	"github.com/liy/goe/plumbing"
)

var Done = errors.New("no more items in iterator")

type CommitIterator struct {
	indegree map[plumbing.Hash]int
	queue   *CommitPrioQueue
	r       *Repository
	Size int
}

func NewCommitIterator(r *Repository, tips []*object.Commit) *CommitIterator {
	itr := &CommitIterator{
		indegree: make(map[plumbing.Hash]int),
		queue:   &CommitPrioQueue{},
		r:       r,
	}

	itr.Size = itr.prepare(tips)

	return itr
}

/*
Traverse function traverses the commit tree and calculate the indegree of each commit.
*/
func (ci *CommitIterator) traverse(c *object.Commit, visited map[plumbing.Hash]bool) {
	visited[c.Hash] = true
	
	for _, ph := range c.Parents {
		pc, _ := ci.r.GetCommit(ph)
		ci.indegree[ph]++

		if !visited[ph] {
			ci.traverse(pc, visited)
		}
	}
}
/*
Prepare the commit iterator by calculating the indegree of each commit. Duplicated tip commits are ignored.
*/
func (ci *CommitIterator) prepare(tips []*object.Commit) int {
	visited := make(map[plumbing.Hash]bool)

	for _, c := range tips {
		if !visited[c.Hash] {
			ci.indegree[c.Hash] = 0
			ci.traverse(c, visited)
		}
	}

	// Removing the duplicated tip commits
	queued := make(map[plumbing.Hash]bool, len(tips))
	for _, t := range tips {
		if ci.indegree[t.Hash] == 0 && !queued[t.Hash] {
			queued[t.Hash] = true
			ci.queue.Enqueue(t)
		}
	}

	return len(visited)
}

func (ci *CommitIterator) Next() (*object.Commit, error) {
	if ci.queue.Size() == 0 {
		return nil, Done
	}

	commit := ci.queue.Dequeue()

	// Try enqueue next commit's parents, if and only if the indegree of the parent is 0
	for _, ph := range commit.Parents {
		ci.indegree[ph]--
		if ci.indegree[ph] == 0 {
			parent, err := ci.r.GetCommit(ph)
			if err != nil {
				return commit, errors.New("cannot get parent commit: " + ph.String())
			}
			
			ci.queue.Enqueue(parent)
		}
	}

	return commit, nil
}
