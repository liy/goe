package goe

import (
	"errors"

	"github.com/liy/goe/object"
	"github.com/liy/goe/plumbing"
	"github.com/liy/goe/utils"
)

var Done = errors.New("no more items in iterator")

type CommitIterator struct {
	visited map[plumbing.Hash]bool
	queue *utils.PrioQueue
	r *Repository
}

func NewCommitIterator(r *Repository, start *object.Commit) *CommitIterator {
	return &CommitIterator{
		visited:  make(map[plumbing.Hash]bool),
		queue: utils.NewPrioQueue(start),
		r: r,
	}
}

func (ci *CommitIterator) Next() (*object.Commit, error) {
	if ci.queue.Size() == 0 {
		return nil, Done
	}

	// try to get the next commit
	commit, ok := (*ci.queue.Dequeue()).(*object.Commit)
	if !ok {
		return nil, errors.New("not a commit object")
	}

	// enqueue next commit's parents
	for _, ph := range commit.Parents {
		if _, exist := ci.visited[ph]; exist {
			continue
		}
		
		ci.visited[ph] = true
		p, err := ci.r.GetCommit(ph)
		if err != nil {
			return commit, errors.New("cannot get parent commit: " + ph.String()) 
		}

		ci.queue.Enqueue(p)
	}

	return commit, nil
}