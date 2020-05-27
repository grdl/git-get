package new

import "testing"

func TestBranchStatusLocal(t *testing.T) {
	tr := NewRepoWithCommit(t)
	tr.NewBranch("branch")

	repo, err := OpenRepo(tr.Path)
	checkFatal(t, err)

	err = repo.LoadStatus()
	checkFatal(t, err)

	if repo.Status.Branches["master"].Upstream != nil {
		t.Errorf("'master' branch should not have an upstream")
	}

	if repo.Status.Branches["branch"].Upstream != nil {
		t.Errorf("'branch' branch should not have an upstream")
	}
}

func TestBranchStatusCloned(t *testing.T) {
	origin := NewRepoWithCommit(t)

	clone := origin.Clone()
	clone.NewBranch("local")

	repo, err := OpenRepo(clone.Path)
	checkFatal(t, err)

	err = repo.LoadStatus()
	checkFatal(t, err)

	if repo.Status.Branches["master"].Upstream == nil {
		t.Errorf("'master' branch should have an upstream")
	}

	if repo.Status.Branches["local"].Upstream != nil {
		t.Errorf("'local' branch should not have an upstream")
	}

}
