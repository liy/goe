package plumbing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsRemote(t *testing.T) {
	assert.False(t, IsRemote("refs/heads/dev"), "refs/heads/dev is not remote")
	assert.False(t, IsRemote("refs/tags/dev"), "refs/tags/dev is not remote")
	assert.False(t, IsRemote("refs/notes/dev"), "refs/notes/dev is not remote")
	assert.True(t, IsRemote("refs/remotes/origin/feat/123"), "refs/remotes/origin/feat/123 is a remote")
	assert.True(t, IsRemote("refs/remotes/origin/dev"), "refs/remotes/notes/dev is a remote")
}

func TestIsBranch(t *testing.T) {
	assert.True(t, IsBranch("refs/heads/dev"), "refs/heads/dev is a branch")
	assert.False(t, IsBranch("refs/tags/dev"), "refs/tags/dev is not branch")
	assert.False(t, IsBranch("refs/notes/dev"), "refs/notes/dev is not branch")
	assert.True(t, IsBranch("refs/remotes/origin/dev"), "refs/remotes/notes/dev is a branch")
	assert.True(t, IsBranch("refs/remotes/origin/feat/123"), "refs/remotes/origin/feat/123 is a branch")
	assert.False(t, IsBranch("refs/remotes/origin/HEAD"), "refs/remotes/origin/HEAD is a branch")
}

func TestIsTag(t *testing.T) {
	assert.False(t, IsTag("refs/heads/dev"), "refs/heads/dev is not a tag")
	assert.True(t, IsTag("refs/tags/dev"), "refs/tags/dev is a tag")
	assert.False(t, IsTag("refs/notes/dev"), "refs/notes/dev is not tag")
	assert.False(t, IsTag("refs/remotes/origin/dev"), "refs/remotes/notes/dev is not a tag")
	assert.False(t, IsTag("refs/remotes/origin/feat/123"), "refs/remotes/origin/feat/123 is not a tag")
	assert.False(t, IsTag("refs/remotes/origin/HEAD"), "refs/remotes/origin/HEAD is not a tag")
}

func TestIsHead(t *testing.T) {
	assert.False(t, IsHead("refs/heads/HEAD"), "refs/heads/HEAD is not a head")
	assert.False(t, IsHead("refs/tags/HEAD"), "refs/tags/HEAD is a head")
	assert.False(t, IsHead("refs/notes/HEAD"), "refs/notes/HEAD is not head")
	assert.False(t, IsHead("refs/remotes/origin/feat/HEAD"), "refs/remotes/origin/feat/123 is not a head")
	assert.True(t, IsHead("refs/remotes/origin/HEAD"), "refs/remotes/origin/HEAD is not a head")
	assert.True(t, IsHead("HEAD"), "HEAD is a head")
}

func TestIsReference(t *testing.T) {
	assert.True(t, ReferenceTarget("ref: refs/heads/master").IsReference(), "is a reference")
	assert.False(t, ReferenceTarget("f2010ee942a47bec0ca7e8f04240968ea5200735").IsReference(), "is not a reference")
	assert.False(t, ReferenceTarget("f2010ee942a47bec0ca7e8f04240968ea5200735").IsReference(), "is not a reference")
}

func TestIsHash(t *testing.T) {
	assert.False(t, ReferenceTarget("ref: refs/heads/master").IsHash(), "is not a hash")
	assert.True(t, ReferenceTarget("f2010ee942a47bec0ca7e8f04240968ea5200735").IsHash(), "is a hash")
}

func TestIsSymbolic(t *testing.T) {
	head := Reference {
		Name: "HEAD",
		Target: "ref: refs/heads/master",
	}
	assert.True(t, head.IsSymbolic(), "HEAD is symbolic")

	master := Reference {
		Name: "refs/heads/master",
		Target: "f2010ee942a47bec0ca7e8f04240968ea5200735",
	}
	assert.False(t, master.IsSymbolic(), "master branch is not symbolic")
}