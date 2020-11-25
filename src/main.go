//+ build js,wasm

package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"syscall/js"

	"github.com/MichaelMure/git-bug/cache"
	"github.com/happybeing/webpack-golang-wasm-async-loader/gobridge"

	// OK FOR CLI if I have gobridge/go.mod containing:
	// module github.com/happybeing/webpack-golang-wasm-async-loader/gobridge
	// go 1.13

	"github.com/MichaelMure/git-bug/cache"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

var global = js.Global()

var fs = memfs.New()

func uploadFile(this js.Value, args []js.Value) (interface{}, error) {
	// if !ready {
	// 	return nil, nil
	// }

	ret := 0
	fullPath := args[0].String()

	array := args[1]
	println("Array byteLength: ", array.Get("byteLength").Int())
	buf := make([]byte, array.Get("byteLength").Int())
	n := js.CopyBytesToGo(buf, array)

	fmt.Println("GO uploading: ", fullPath, n)

	dst, err := fs.Create(fullPath)
	if err != nil {
		return nil, err
	}

	_, err = dst.Write(buf)
	if err != nil {
		return nil, err
	}

	if err = dst.Close(); err != nil {
		return nil, err
	}

	return ret, nil
}

func listFiles(this js.Value, args []js.Value) (interface{}, error) {
	listing, err := fs.ReadDir("testrepo")
	// listing, err := fs.ReadDir(args[0].String())
	if err != nil {
		log.Fatal(err)
	}

	println("Listing files:")
	for _, f := range listing {
		fmt.Println(f.Name())
	}

	return 0, err
}

//// Manage repositories
//
// This approach is just for the proof-of-concept. When using a p2p filesystem,
// we probably want to clone directly to that rather than a Storer, but will
// need to clone into a filesystem in memory rather than an in memory Storer.
//
// For now we keep a map of cloned URIs to Repository objects.
//

type RepoEntry struct {
	host      string
	path      string
	url       string
	gogitRepo *gogit.Repository
}

var repositories = make(map[string]*RepoEntry, 0)

func getRepositoryList(this js.Value, args []js.Value) (interface{}, error) {
	retRepos := make([]interface{}, len(repositories))

	repoIndex := 0
	for path, entry := range repositories {
		cfg, err := entry.gogitRepo.Config()
		if err != nil {
			return nil, nil
		}
		repo := make(map[string]interface{}, 0)
		repo["path"] = path
		repo["host"] = entry.host
		repo["path"] = entry.path
		repo["url"] = entry.url
		repo["author"] = cfg.Author.Name
		repo["author-email"] = cfg.Author.Email
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
func cloneRepository(this js.Value, args []js.Value) (interface{}, error) {
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

	storage := memory.NewStorage()
	var err error
	go func() {
		r, err := git.Clone(storage, fs, &git.CloneOptions{URL: url})
		if err != nil {
			println("git.Clone() failed: ", err.Error())
			callback.Invoke(err.Error())
		} else {
			entry := RepoEntry{
				host,
				path,
				url,
				r,
			}
			repositories[path] = &entry
			callback.Invoke(js.Null())
		}
	}()

	return nil, err
}

// NOTE: must disable CORS in the browser using a browser plugin
func testGitClone(this js.Value, args []js.Value) (interface{}, error) {
	message := ""
	host := ""
	path := ""
	url := ""

	// A small test GitLab repo:

	// Testing ways around CORS...

	// gitlab.com
	// WORKS with gitlab.com if browser CORS disabled using a plugin
	message = "https clone of gitlab repo"
	host = "https://gitlab.com/"
	path = "saeedareffard1377666/testproject2.git"

	// github.com
	// WORKS using local-cors-proxy, only for a specific service:
	// yarn global add local-cors-proxy
	// lcp --proxyUrl https://github.com
	host = "http://localhost:8010/proxy/"
	path = "happybeing/p2p-git-portal-poc.git"
	url = host + path

	// WORKS with isomorphic-git proxy service:
	// url = "https://cors.isomorphic-git.org/github.com/happybeing/p2p-git-portal-poc.git"

	// FAILS without proxy (including with CORS disabled in browser)
	// url = "https://github.com/happybeing/p2p-git-portal-poc.git"
	// token := "<replace with PAT>"

	println(message, url)
	storage := memory.NewStorage()
	// Note: can't wait on this or go will exit due to a deadlock (as Clone() uses http.Get())
	go func() {
		r, err := git.Clone(storage, fs, &git.CloneOptions{
			// Auth: &http.BasicAuth{
			// 	Username: "noname", // This can be anything except an empty string
			// 	Password: token,
			// },
			URL: url,
		})
		if err != nil {
			println("git.Clone() failed: ", err.Error())
		} else {
			println("Clone complete...")
			entry := RepoEntry{
				host,
				path,
				url,
				r,
			}
			repositories[path] = &entry

			// Retrieve the branch pointed by HEAD
			ref, err2 := r.Head()
			if err2 != nil {
				println("r.Head() failed: ", err2.Error())
			} else {
				println("Retrieved head, ref: ", ref)
			}
		}
	}()

	return nil, nil
}

//// Repository information
//
// showcase example: https://github.com/go-git/go-git/blob/master/_examples/showcase/main.go

// args[]:
//	[0] path of a cloned repository

func listHeadCommits(this js.Value, args []js.Value) (interface{}, error) {
	path := args[0].String()

	println("arg path: ", path)
	// path = "happybeing/p2p-git-portal-poc.git"

	entry, found := repositories[path]
	if !found {
		println("Unknown repository: ", path)
		return nil, nil
	}
	r := entry.gogitRepo

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
	commitIter, err := r.Log(&git.LogOptions{From: commit.Hash})
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
//	[1] first - index of first commit starting at zero
//  [2] last - index of last commit (inclusive, so to return just the first commit first=last=0)

func getHeadCommitsRange(this js.Value, args []js.Value) (interface{}, error) {
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

	entry, found := repositories[path]
	if !found {
		println("Unknown repository: ", path)
		return nil, nil
	}
	r := entry.gogitRepo

	// Latest commit on current branch
	ref, err := r.Head()
	if err != nil {
		return nil, err
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	// List commits
	commitIter, err := r.Log(&git.LogOptions{From: commit.Hash})
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

//// Test syscall/js Go/Wasm types

func testTypes(this js.Value, args []js.Value) (interface{}, error) {
	person := make(map[string]interface{}, 0)

	person["name"] = "Alice"
	person["age"] = 35
	person["height"] = 167.64

	child := make(map[string]interface{}, 0)
	child["name"] = "Peter"
	child["age"] = 10
	person["child"] = child

	return person, nil
}

////// git-bug

func testGitBug(this js.Value, args []js.Value) (interface{}, error) {
	cache, err := cache.NewMultiRepoCache()
	if err != nil {
		return nil, err
	}
	println("testGitBug() created multi-repo cache!", cache)
	return nil, nil
}

////// Go/wasm initialisation

var ready = false

func main() {
	c := make(chan struct{}, 0)

	gobridge.RegisterCallback("uploadFile", uploadFile)
	gobridge.RegisterCallback("listFiles", listFiles)
	gobridge.RegisterCallback("listHeadCommits", listHeadCommits)
	gobridge.RegisterCallback("testTypes", testTypes)
	gobridge.RegisterCallback("testGitClone", testGitClone)

	gobridge.RegisterCallback("cloneRepository", cloneRepository)
	gobridge.RegisterCallback("getRepositoryList", getRepositoryList)
	gobridge.RegisterCallback("getHeadCommitsRange", getHeadCommitsRange)

	gobridge.RegisterCallback("testGitBug", testGitBug)
	ready = true
	println("Web Assembly is ready")
	<-c // Makes the Go process wait until we want it to end
}
