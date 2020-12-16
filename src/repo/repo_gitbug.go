//+ build js,wasm

// Package repo manages a set of git-bug gogit.Repository each with go-billy fs
package repo

import (
	"path/filepath"
	"syscall/js"

	bugcache "github.com/MichaelMure/git-bug/cache"
	gogit "github.com/go-git/go-git/v5"

	"github.com/MichaelMure/git-bug/repository"
	"github.com/go-git/go-git/v5/plumbing/cache"
	storagefs "github.com/go-git/go-git/v5/storage/filesystem"
)

var GitbugCache = bugcache.NewMultiRepoCache()
var GitbugIds = make(map[string]*bugcache.IdentityCache, 0)

// Test using git-bug API

func SetupGitbugCache(this js.Value, args []js.Value) (interface{}, error) {
	println("SetupGitbugCache()...")
	if GitbugCache == nil {
		println("repo_gitbug.go: GitbugCache was not created")
	}

	entryPrimary, err := NewGitbugRepoCache("primary-repo")
	if err != nil {
		return nil, err
	}
	entryAnother, err := NewGitbugRepoCache("another-repo")

	rene, _ := NewTestGitbugIdentity(entryPrimary, "René Descartes", "rene@descartes.fr")
	mark, _ := NewTestGitbugIdentity(entryPrimary, "Mark Hughes", "gitbug@markhughes.com")
	AddSampleBugs(entryPrimary, mark)

	err = entryAnother.GitbugRC.SetUserIdentity(rene)
	if err != nil {
		return nil, err
	}

	// TODO add code to create test commits (probably as a separate function)

	// Test list bugs from cache
	GetBugsForRepo("primary-repo")

	// Test open repo from filesystem and list bugs
	GitbugCache.Close()
	GitbugCache = bugcache.NewMultiRepoCache()

	return GetBugsForRepo("primary-repo")
}

// args []js.Value:
// [0] - path to worktree where repository will be created
// [1] - optional callback(error) called on completion with a null value on success, or an error message
func OpenRepository(this js.Value, args []js.Value) (interface{}, error) {
	println("OpenRepository()...")
	path := args[0].String()

	var callback js.Value
	if len(args) > 1 {
		callback = args[1]
	}

	// worktreeFs, err := Filesystem.Chroot(path)
	// if err != nil {
	// 	return nil, err
	// }

	// dotGitFs, err := Filesystem.Chroot(filepath.Join(path, ".git"))
	// if err != nil {
	// 	return nil, err
	// }

	// storage := storagefs.NewStorage(dotGitFs, cache.NewObjectLRUDefault())
	// gogitRepo, err := gogit.Init(storage, worktreeFs)
	// if err != nil {
	// 	println("gogit.Init() failed: ", err.Error())
	// 	return nil, err
	// }
	gogitRepo, repoCache, err := OpenRepo(path)
	if err != nil {
		println("OpenRepo() failed: ", err.Error())
		callback.Invoke(err.Error())
	} else {
		entry := Entry{
			"none",
			path,
			"none",
			gogitRepo,
			repoCache,
		}

		// bilbo, _ := NewTestGitbugIdentity(&entry, "Bilbo Baggins", "bilbo@the.vale")
		// AddSampleBugs(&entry, bilbo)
		AllRepositories[path] = &entry
		callback.Invoke(js.Null())
	}

	return nil, nil
}

// args []js.Value:
// [0] - path to worktree where repository will be created
func NewRepository(this js.Value, args []js.Value) (interface{}, error) {
	println("NewRepository()...")
	path := args[0].String()

	worktreeFs, err := Filesystem.Chroot(path)
	if err != nil {
		return nil, err
	}

	dotGitFs, err := Filesystem.Chroot(filepath.Join(path, ".git"))
	if err != nil {
		return nil, err
	}

	storage := storagefs.NewStorage(dotGitFs, cache.NewObjectLRUDefault())
	gogitRepo, err := gogit.Init(storage, worktreeFs)
	if err != nil {
		println("gogit.Init() failed: ", err.Error())
		return nil, err
	}
	_, repoCache, err := OpenRepo(path)
	if err != nil {
		println("OpenRepo() failed: ", err.Error())
		return nil, err
	}

	entry := Entry{
		"none",
		path,
		"none",
		gogitRepo,
		repoCache,
	}

	bilbo, _ := NewTestGitbugIdentity(&entry, "Bilbo Baggins", "bilbo@the.vale")
	AddSampleBugs(&entry, bilbo)
	AllRepositories[path] = &entry

	return nil, nil
}

func AddSampleBugs(entry *Entry, user *bugcache.IdentityCache) {
	println("AddSampleBugs() to ", entry.Path)
	repo := entry.GitbugRC
	err := repo.SetUserIdentity(user)
	if err != nil {
		return
	}

	// Add a few bugs
	_, _, err = repo.NewBug("This bug is a triumph!", "Lots to say\nand\nno\ntime to\nsay\nit!")
	_, _, err = repo.NewBug("One bug is good but two is better", "Nothing to see here though.")
	bug3, _, err := repo.NewBug("Style the HTML/CSS to suit a git portal", "The app runs and uses default Svelte HTML/CSS styling, but initially it should mirror github to highlight its aim to replicate key features of a git portal like github. Later it can differentiate itself, but initially I think it will help to demonstrate the purpose of a proof-of-concept, what is working at any stage and where work is still needed.\n	As a starting point, we might look at using parts of the git-bug web UI, or build something from scratch for this Svelte framework.")

	bug3.AddComment("First comment on a bug!")

}

func NewGitbugRepoCache(path string) (*Entry, error) {
	repo, err := repository.InitFsGoGitRepo(path, Filesystem)
	if err != nil {
		return nil, err
	}
	repoCache, err := GitbugCache.RegisterRepository(path, repo)
	if err != nil {
		return nil, err
	}

	entry := Entry{
		"none",
		path,
		"none",
		repo.GetGitRepo(),
		repoCache,
	}
	AllRepositories[path] = &entry

	return &entry, nil
}

func NewTestGitbugIdentity(entry *Entry, name string, email string) (*bugcache.IdentityCache, error) {
	newId, err := entry.GitbugRC.NewIdentity(name, email)
	// TODO push id to other repos (first add as remotes to current repo)
	if err != nil {
		return nil, err
	}

	GitbugIds[name] = newId
	return newId, nil
}

func OpenRepo(path string) (*gogit.Repository, *bugcache.RepoCache, error) {
	existingEntry := AllRepositories[path]
	if existingEntry != nil {
		return existingEntry.GogitRepo, existingEntry.GitbugRC, nil
	}

	rc, _ := GitbugCache.ResolveRepo(path)
	repo, err := repository.OpenFsGoGitRepo(path, nil, Filesystem)
	if err != nil {
		return nil, nil, err
	}

	if rc == nil {
		rc, err = GitbugCache.RegisterRepository(path, repo)
		if err != nil {
			return nil, nil, err
		}
	}

	return repo.GetGitRepo(), rc, nil
}

func GetBugsForRepo(path string) ([]interface{}, error) {
	_, rc, err := OpenRepo(path)
	if err == nil {
		return listBugsForRepoCache(rc, nil)
	}
	return nil, err
}

func listBugsForRepoCache(repoCache *bugcache.RepoCache, err error) ([]interface{}, error) {
	println("listBugsForRepoCache()...")
	if err != nil || repoCache == nil {
		return nil, err
	}

	bugs := make([]interface{}, len(repoCache.AllBugsIds()))
	for bugIndex, id := range repoCache.AllBugsIds() {
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
			bugs[bugIndex] = bug
		}
	}

	return bugs, nil
}
