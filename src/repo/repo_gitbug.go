//+ build js,wasm

// Package repo manages a set of git-bug gogit.Repository each with go-billy fs
package repo

import (
	"syscall/js"

	"github.com/MichaelMure/git-bug/cache"
	"github.com/MichaelMure/git-bug/repository"
)

var GitbugCache = cache.NewMultiRepoCache()
var GitbugIds = make(map[string]*cache.IdentityCache, 0)

// Test using git-bug API

func SetupGitbugCache(this js.Value, args []js.Value) (interface{}, error) {
	println("SetupGitbugCache()...")
	if GitbugCache == nil {
		println("repo_gitbug.go: GitbugCache was not created")
	}

	primaryRepo, err := createTestGitbugRepoCache("primary-repo")
	if err != nil {
		return nil, err
	}

	anotherRepo, err := createTestGitbugRepoCache("another-repo")
	rene, _ := createTestGitbugIdentity(primaryRepo, "René Descartes", "rene@descartes.fr")
	mark, _ := createTestGitbugIdentity(anotherRepo, "Mark Hughes", "gitbug@markhughes.com")

	err = primaryRepo.SetUserIdentity(rene)
	if err != nil {
		return nil, err
	}

	// Add a few bugs
	_, _, err = primaryRepo.NewBug("This bug is a triumph!", "Lots to say\nand\nno\ntime to\nsay\nit!")
	_, _, err = primaryRepo.NewBug("One bug is good but two is better", "Nothing to see here though.")
	bug3, _, err := primaryRepo.NewBug("Style the HTML/CSS to suit a git portal", "The app runs and uses default Svelte HTML/CSS styling, but initially it should mirror github to highlight its aim to replicate key features of a git portal like github. Later it can differentiate itself, but initially I think it will help to demonstrate the purpose of a proof-of-concept, what is working at any stage and where work is still needed.\n	As a starting point, we might look at using parts of the git-bug web UI, or build something from scratch for this Svelte framework.")

	bug3.AddComment("First comment on a bug!")

	err = anotherRepo.SetUserIdentity(mark)
	if err != nil {
		return nil, err
	}

	// TODO add code to create test commits (probably as a separate function)

	// Test list bugs from cache
	listBugsForRepo("primary-repo")

	// Test open repo from filesystem and list bugs
	GitbugCache.Close()
	GitbugCache = cache.NewMultiRepoCache()

	return listBugsForRepo("primary-repo")
}

func createTestGitbugRepoCache(path string) (*cache.RepoCache, error) {
	repo, err := repository.InitFsGoGitRepo(path, Filesystem)
	if err != nil {
		return nil, err
	}

	repoCache, err := GitbugCache.RegisterRepository(path, repo)
	if err != nil {
		return nil, err
	}

	return repoCache, nil
}

func createTestGitbugIdentity(rc *cache.RepoCache, name string, email string) (*cache.IdentityCache, error) {
	newId, err := rc.NewIdentity(name, email)
	// TODO push id to other repos (first add as remotes to current repo)
	if err != nil {
		return nil, err
	}

	GitbugIds[name] = newId
	return newId, nil
}

func openRepo(path string) (*cache.RepoCache, error) {
	rc, _ := GitbugCache.ResolveRepo(path)
	if rc != nil {
		return rc, nil
	}

	repo, err := repository.OpenFsGoGitRepo(path, nil, Filesystem)
	if err != nil {
		return nil, err
	}

	rc, err = GitbugCache.RegisterRepository(path, repo)
	if err != nil {
		return nil, err
	}

	return rc, nil
}

func listBugsForRepo(path string) (map[string]interface{}, error) {
	return listBugsForRepoCache(openRepo(path))
}

func listBugsForRepoCache(repoCache *cache.RepoCache, err error) (map[string]interface{}, error) {
	println("listBugsForRepoCache()...")
	if err != nil || repoCache == nil {
		return nil, err
	}

	bugs := make(map[string]interface{}, 0)
	for _, id := range repoCache.AllBugsIds() {
		// println("bugId:", (string)(id))
		bugCache, err := repoCache.ResolveBug(id)
		if err != nil {
			println("err:", err.Error())
		} else {
			bug := make(map[string]interface{}, 0)
			snap := bugCache.Snapshot()
			println("title: ", snap.Title)
			bug["title"] = snap.Title
			bug["id"] = snap.Id().Human()
			bug["author-name"] = snap.Author.DisplayName()
			bug["created"] = snap.CreateTime.String()
			labels := make([]string, len(snap.Labels))
			for i, l := range snap.Labels {
				labels[i] = l.String()
			}
		}
	}

	return bugs, nil
}
