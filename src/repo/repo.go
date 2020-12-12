//+ build js,wasm

// Package repo manages a set of git-bug gogit.Repository each with go-billy fs
package repo // TODO maybe rename to 'api' or something not 'repo'

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"syscall/js"

	// "github.com/go-git/go-git/v5"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"

	bugcache "github.com/MichaelMure/git-bug/cache"
	// "github.com/go-git/go-git/storage/memory"
	gogit "github.com/go-git/go-git/v5"
	storagefs "github.com/go-git/go-git/v5/storage/filesystem"
)

// Entry holds metadata and for a gogit.Repository
type Entry struct {
	Host      string
	Path      string
	URL       string
	GogitRepo *gogit.Repository
	GitbugRC  *bugcache.RepoCache
}

var Filesystem = memfs.New()
var AllRepositories = make(map[string]*Entry, 0)

func GetRepositoryList(this js.Value, args []js.Value) (interface{}, error) {
	retRepos := make([]interface{}, len(AllRepositories))

	repoIndex := 0
	for path, entry := range AllRepositories {
		// cfg, err := entry.GogitRepo.Config()
		// if err != nil {
		// 	return nil, nil
		// }
		repo := make(map[string]interface{}, 0)
		repo["path"] = path
		repo["host"] = entry.Host
		repo["path"] = entry.Path
		repo["url"] = entry.URL
		// repo["author"] = cfg.Author.Name
		// repo["author-email"] = cfg.Author.Email
		retRepos[repoIndex] = repo
		repoIndex += 1
	}

	return retRepos, nil
}

// Clone a git repository
//
// This works under the following conditions:
// - cloning a gitlab URL such as: https://gitlab.com/saeedareffard1377666/testproject2.git
//   if you disable CORS in the browser (eg in Chromium with the Allow CORS plugin)
// - cloning from gitlab.com and github.com works using a proxy. You can
//   either use a service as in https://cors.isomorphic-git.org/github.com/happybeing/p2p-git-portal-poc.git
//   or run a local proxy. cors-buster WORKS locally:
// 		git clone https://github.com/wmhilton/cors-buster
// 		cd cors-buster && yarn && yarn start
//		# Example "http://localhost:3000/github.com/happybeing/p2p-git-portal-poc.git"
//
// Notes:
// - for large repositories the download can take a long time and in the browser
// console it may look as if nothing is happening and there is no response to a POST requesting
// the packfiles to be downloaded.
//
// - you must do the git.Clone() in separate goroutine and NOT use a WaitGroup to wait
// for the goroutine to complete, or you will cause a deadlock and Go will exit with an
// error. This is due to http requests being issued in the call.
// More: https://golang.org/pkg/syscall/js/#FuncOf
//
// - we can't clone a local file:/// URI, but can upload a local repo with drag and drop
//   in the browser, using Chromium at least.
//
// TODO: Issues to open
// - Implement a way to show progress during a git operation (I think go-git supports this).
// - CORS issues might be handled better but this is not a high priority
//
// args []js.Value:
// [0] - host URI (e.g. 'https://github.com' or 'https://github.com/')
// [1] - repository path without leading '/' (e.g. 'happybeing/p2p-git-portal-poc')
// [2] - optional proxied URI (e.g. 'https://localhost:8171/happybeing/p2p-git-portal-poc')
// [3] - optional callback(error) called on completion with a null value on success, or an error message

func CloneRepository(this js.Value, args []js.Value) (interface{}, error) {
	host := args[0].String()
	if host[len(host)-1] != '/' {
		host = host + "/"
	}

	path := args[1].String()
	var proxiedURI = ""
	if len(args) > 2 {
		proxiedURI = args[2].String()
	}

	var callback js.Value
	if len(args) > 3 {
		callback = args[3]
	}

	url := host + path
	if len(proxiedURI) > 0 {
		url = proxiedURI
	}

	// println("cloneRepository() with url:", url)
	// println("host:", host)
	// println("path:", path)
	// println("proxiedURI:", proxiedURI)
	dotGitFs, err := Filesystem.Chroot(filepath.Join(path, ".git"))
	if err != nil {
		return nil, err
	}

	storage := storagefs.NewStorage(dotGitFs, cache.NewObjectLRUDefault())
	go func() {
		gogitRepo, err := gogit.Clone(storage, Filesystem, &gogit.CloneOptions{URL: url})
		if err != nil {
			// if true {
			println("gogit.Clone() failed: ", err.Error())
			callback.Invoke(err.Error())
		} else {
			repoCache, err := OpenRepo(path)
			if err != nil {
				callback.Invoke(err.Error())
			}
			entry := Entry{
				host,
				path,
				url,
				gogitRepo,
				repoCache,
			}
			bilbo, _ := NewTestGitbugIdentity(&entry, "Bilbo Baggins", "bilbo@the.vale")
			AddSampleBugs(&entry, bilbo)
			AllRepositories[path] = &entry
			callback.Invoke(js.Null())
		}
	}()

	return nil, nil
}

// NOTE: must disable CORS in the browser using a browser plugin
func GitCloneTest(this js.Value, args []js.Value) (interface{}, error) {
	// message := ""
	// host := ""
	// path := ""
	// url := ""

	// // A small test GitLab repo:

	// // Testing ways around CORS...

	// // gitlab.com
	// // WORKS with gitlab.com if browser CORS disabled using a plugin
	// message = "https clone of gitlab repo"
	// host = "https://gitlab.com/"
	// path = "saeedareffard1377666/testproject2.git"

	// // github.com
	// // WORKS using local-cors-proxy, only for a specific service:
	// // yarn global add local-cors-proxy
	// // lcp --proxyUrl https://github.com
	// host = "http://localhost:8010/proxy/"
	// path = "happybeing/p2p-git-portal-poc.git"
	// url = host + path

	// // WORKS with isomorphic-git proxy service:
	// // url = "https://cors.isomorphic-git.org/github.com/happybeing/p2p-git-portal-poc.git"

	// // FAILS without proxy (including with CORS disabled in browser)
	// // url = "https://github.com/happybeing/p2p-git-portal-poc.git"
	// // token := "<replace with PAT>"

	// println(message, url)
	// storage := memory.NewStorage()
	// // Note: can't wait on this or go will exit due to a deadlock (as Clone() uses http.Get())
	// go func() {
	// 	r, err := gogit.Clone(storage, Filesystem, &gogit.CloneOptions{
	// 		// Auth: &http.BasicAuth{
	// 		// 	Username: "noname", // This can be anything except an empty string
	// 		// 	Password: token,
	// 		// },
	// 		URL: url,
	// 	})
	// 	if err != nil {
	// 		println("gogit.Clone() failed: ", err.Error())
	// 	} else {
	// 		println("Clone complete...")
	// 		entry := Entry{
	// 			host,
	// 			path,
	// 			url,
	// 			r,
	// 		}
	// 		AllRepositories[path] = &entry

	// 		// Retrieve the branch pointed by HEAD
	// 		ref, err2 := r.Head()
	// 		if err2 != nil {
	// 			println("r.Head() failed: ", err2.Error())
	// 		} else {
	// 			println("Retrieved head, ref: ", ref)
	// 		}
	// 	}
	// }()

	return nil, nil
}

//// Repository information
//
// showcase example: https://github.com/go-git/go-git/blob/master/_examples/showcase/main.go

// args[]:
//	[0] path of a cloned repository

func ListHeadCommits(this js.Value, args []js.Value) (interface{}, error) {
	path := args[0].String()

	println("arg path: ", path)
	// path = "happybeing/p2p-git-portal-poc.git"

	entry, found := AllRepositories[path]
	if !found {
		println("Unknown repository: ", path)
		return nil, nil
	}
	r := entry.GogitRepo
	// Latest commit on current branch
	ref, err := r.Head()
	if err != nil {
		return nil, err
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}
	fmt.Println("Latest commit:", commit)

	// List commits
	commitIter, err := r.Log(&gogit.LogOptions{From: commit.Hash})
	if err != nil {
		return nil, err
	}

	err = commitIter.ForEach(func(c *object.Commit) error {
		hash := c.Hash.String()
		line := strings.Split(c.Message, "\n")
		fmt.Println(hash[:7], line[0])

		return nil
	})

	return nil, nil
}

// args[]:
//	[0] path - local repo path of a cloned repository (later will be a local identifier)
func GetIssuesForRepo(this js.Value, args []js.Value) (interface{}, error) {
	path := args[0].String()
	// entry, found := AllRepositories[path]
	// if !found {
	// 	println("Unknown repository: ", path)
	// 	return nil, nil
	// }
	// r := entry.GogitRepo

	return GetBugsForRepo(path)
}

// args[]:
//	[0] path - local repo path of a cloned repository (later will be a local identifier)
//	[1] first - index of first commit starting at zero
//  [2] last - index of last commit (inclusive, so to return just the first commit first=last=0)

func GetHeadCommitsRange(this js.Value, args []js.Value) (interface{}, error) {
	path := args[0].String()
	first := args[1].Int()
	last := args[2].Int()

	var err error
	if last < first {
		err = errors.New("Invalid range, 'last' must be at least 'first'")
	} else if first < 0 {
		err = errors.New("Invalid range, 'first' must not be less than 0")
	}
	if err != nil {
		return nil, err
	}

	println("arg path: ", path)
	// path = "saeedareffard1377666/testproject2.git"

	entry, found := AllRepositories[path]
	if !found {
		println("Unknown repository: ", path)
		return nil, nil
	}
	r := entry.GogitRepo

	// Latest commit on current branch
	ref, err := r.Head()
	if ref == nil || err != nil {
		return nil, err
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	// List commits
	commitIter, err := r.Log(&gogit.LogOptions{From: commit.Hash})
	if err != nil {
		return nil, err
	}

	commits := make([]interface{}, last-first+1)
	commitIndex := 0
	totalCommits := 0

	err = commitIter.ForEach(func(c *object.Commit) error {
		totalCommits += 1
		if commitIndex >= first && commitIndex <= last {
			commit := make(map[string]interface{}, 0)
			commit["hash"] = c.Hash.String()
			commit["message"] = c.Message

			commits[commitIndex] = commit
			commitIndex += 1
		}
		return nil
	})

	retCommits := make(map[string]interface{}, 0)
	retCommits["length"] = commitIndex
	retCommits["totalCommits"] = totalCommits
	retCommits["commits"] = commits[0:commitIndex]

	return retCommits, nil
}
