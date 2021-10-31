package git

import (
	"testing"

	"github.com/liy/goe/plumbing"
	"github.com/liy/goe/tests"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	repo *Repository
}

func (suite *Suite) SetupTest() {
	repo, err := Open("/topo-sort/", tests.Embeded{})
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.repo = repo
}

func (suite *Suite) TestHead() {
	ref, err := suite.repo.HEAD()
	if err != nil {
		suite.T().Fatal(err)
	}

	suite.Equal("refs/heads/dev", ref.Target.ReferenceName(), "HEAD is pointed to dev branch")
}

func (suite *Suite) TestTryPeel() {
	ref, err := suite.repo.HEAD()
	if err != nil {
		suite.T().Fatal(err)
	}

	hash := suite.repo.Peel(ref)
	suite.Equal("f2010ee942a47bec0ca7e8f04240968ea5200735", hash.String(), "HEAD pointed to dev branch tip commit")
}

func (suite *Suite) TestGetReferences() {
	refs := suite.repo.GetReferences()
	
	tests.ToMatchSnapshot(suite.T(), refs)
}

func (suite *Suite) TestGetCommit() {
	c, err := suite.repo.GetCommit(plumbing.ToHash("4854f1ce5137086767e71b4d3010db28bcd09c49"))
	if err != nil {
		suite.T().Fatal(err)
	}

	suite.True(len(c.Parents) == 1, "only has 1 parent")
	suite.Equal("c53c4c18e245d880899405c07eb4d01b735b72ad", c.Parents[0].String(), "has correct parent commit")
	suite.Equal("Fix package keywords.\n", c.Message, "commit has correct message")
	tests.ToMatchSnapshot(suite.T(), c)
}

func (suite *Suite) TestGetCommits() {
	cs, err := suite.repo.GetCommits(plumbing.ToHash("c91773cc3da3c2c3c954626a0b6a44c3ac9e3e92"))
	if err != nil {
		suite.T().Fatal(err)
	}
	
	tests.ToMatchSnapshot(suite.T(), cs)
}

func (suite *Suite) TestGetAnnotatedTag() {
	tag, err := suite.repo.GetAnnotatedTag(plumbing.ToHash("ca0a44b6eddd79547d1ad8bc94be987489edde2a"))
	if err != nil {
		suite.T().Fatal(err)
	}

	tests.ToMatchSnapshot(suite.T(), tag)
}

func (suite *Suite) TestReadObject() {
	obj, err := suite.repo.ReadObject(plumbing.ToHash("4854f1ce5137086767e71b4d3010db28bcd09c49"))
	if err != nil {
		suite.T().Fatal(err)
	}

	tests.ToMatchSnapshot(suite.T(), obj)
}

func TestSuite(t *testing.T) {
    suite.Run(t, new(Suite))
}