package api39

import (
	"fmt"
	"github.com/libgit2/git2go/v31"
)

func fastForward(repo *git.Repository, branch *git.Branch) error {
	// Checkout HEAD
	var ff_checkout_opts git.CheckoutOptions
	ff_checkout_opts.Strategy = git.CheckoutSafe
	if err := repo.CheckoutHead(&ff_checkout_opts); err != nil {
		return err
	}

	// Get a reference to HEAD
	target, err := repo.Head()
	if err != nil {
		return err
	}
	defer target.Free()

	// Set the HEAD reference to the branch's HEAD
	newtarget, err := target.SetTarget(branch.Reference.Target(), "Fast-forwarding")
	if err != nil {
		return err
	}
	defer newtarget.Free()

	// Lookup the commit for the new HEAD reference
	commit, err := repo.LookupCommit(newtarget.Target())
	if err != nil {
		return err
	}

	// Reset the repo to the new HEAD
	return repo.ResetToCommit(commit, git.ResetHard, &ff_checkout_opts)
}

func fetch(repo *git.Repository, remote string) error {
	origin, err := repo.Remotes.Lookup(remote)
	if err != nil {
		return err
	}
	if err = origin.Fetch([]string{}, &git.FetchOptions{}, ""); err != nil {
		return err
	}
	return nil
}

func analyzeForMerge(repo *git.Repository, branch *git.Branch) (git.MergeAnalysis, error) {
	commit, err := repo.AnnotatedCommitFromRef(branch.Reference)
	if err != nil {
		return git.MergeAnalysisNone, err
	}
	defer commit.Free()

	analysis, _, err := repo.MergeAnalysis([]*git.AnnotatedCommit{commit})
	if err != nil {
		return analysis, err
	}

	return analysis, nil
}

func UpdateGitRepo(url, path, branchname string) error {
	if repo, err := git.OpenRepository(path); err != nil {
		_, err = git.Clone(url, path, &git.CloneOptions{})
		return err
	} else {
		defer repo.Free()

		if err := fetch(repo, "origin"); err != nil {
			return err
		}

		branch, err := repo.LookupBranch(branchname, git.BranchAll)
		if err != nil {
			return err
		}

		analysis, err := analyzeForMerge(repo, branch)
		if err != nil {
			return err
		}

		if analysis&git.MergeAnalysisFastForward != 0 {
			return fastForward(repo, branch)
		} else if analysis&git.MergeAnalysisUpToDate != 0 {
			return nil
		} else {
			return fmt.Errorf("cannot fastforward %s branch", branchname)
		}
	}
}
