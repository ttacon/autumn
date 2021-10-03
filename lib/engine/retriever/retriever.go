package retriever

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/ttacon/autumn/lib/config"
)

type FrameworkRetriever interface {
	Get(frmwrk config.FrameworkGetter) error
}

func NewFrameworkRetriever(
	c config.Config,
	roots config.ConfigLoadRoots,
) (FrameworkRetriever, error) {

	// NOTE(ttacon): we're returning an error to reserve the ability to do
	// config validation at a future time.

	return &frameworkRetriever{
		conf: c,
	}, nil
}

type frameworkRetriever struct {
	conf config.Config
}

func (f *frameworkRetriever) Get(frmwrk config.FrameworkGetter) error {
	// NOTE(ttacon): this entire section on protocol joining/creation/handling
	// needs to be cleaned up.
	frameworkURL := frmwrk.GetFramework()
	if len(frameworkURL) == 0 {
		return nil
	}

	slugSafe := strings.ReplaceAll(frameworkURL, "/", "__")

	// This code assumes that the .autumn directory exists.
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	dir := filepath.Join(
		cwd,
		".autumn",
		"frameworks",
		slugSafe,
	)

	r, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL: frmwrk.GetProtocol() + frameworkURL,
	})
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	checkoutOpts := &git.CheckoutOptions{}

	version := frmwrk.GetVersion()
	if len(version) == 0 {
		checkoutOpts.Branch = plumbing.ReferenceName("refs/heads/main")
	} else if strings.HasPrefix(version, "v") {
		checkoutOpts.Branch = plumbing.ReferenceName("refs/tags/" + version)
	} else {
		checkoutOpts.Hash = plumbing.NewHash(version)
	}

	if err = w.Checkout(checkoutOpts); err != nil {
		return err
	}

	return nil
}
