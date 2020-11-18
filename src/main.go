//+ build js,wasm

package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"syscall/js"

	"github.com/happybeing/webpack-golang-wasm-async-loader/gobridge"

	// OK FOR CLI if I have gobridge/go.mod containing:
	// module github.com/happybeing/webpack-golang-wasm-async-loader/gobridge
	// go 1.13
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

var repositories = map[string](*gogit.Repository){}

func listRepositories(this js.Value, args []js.Value) (interface{}, error) {
	println("Repositories:")
	for uri, r := range repositories {
		fmt.Println(uri)
		cfg, err := r.Config()
		if err != nil {
			return nil, nil
		}

		fmt.Println("Author:", cfg.Author.Name, cfg.Author.Email)

		// if f.IsDir() {
		// 	fmt.Println(f.Name())
		// 	descriptionPath := f.Name()
		// 	file, err := fs.Open(descriptionPath)
		// 	if err != nil {
		// 		fmt.Println("can't open ", descriptionPath)
		// 	} else {
		// 		description := make([]byte, 2048)
		// 		_, _ = file.Read(description)
		// 		fmt.Println(f.Name, '-', description)
		// 	}
		// }
	}

	return 0, nil
}

// Clone a git repository
//
// This works under the following conditions:
// - cloning a gitlab URL such as: https://gitlab.com/saeedareffard1377666/testproject2.git
//   if you disable CORS in the browser (eg in Chromium with the Allow CORS plugin)
// - cloning from gitlab.com and github.com works using a proxy. You can
//   either use a service as in https://cors.isomorphic-git.org/github.com/happybeing/p2p-git-portal-poc.git
//   or run a local proxy. This one only works for a single repo domain at a time:
// 		yarn global add local-cors-proxy
// 		lcp --proxyUrl https://github.com
//	 	# Then access as https://localhost:8171/github.com/happybeing/p2p-git-portal-poc.git
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
// - CORS issues might be handled better, such as by running a proxy in the browser using
//   JavaScript but this is not a high priority.
//
func gitClone(this js.Value, args []js.Value) (interface{}, error) {
	url := ""
	message := ""
	println(message, url)

	url = "https://gitlab.com/weblate/libvirt"

	storage := memory.NewStorage()
	go func() {
		r, err := git.Clone(storage, fs, &git.CloneOptions{URL: url})
		if err != nil {
			println("git.Clone() failed: ", err.Error())
		} else {
			println("Clone complete...")
			repositories[url] = r
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

// NOTE: must disable CORS in the browser using a browser plugin
func testGitClone(this js.Value, args []js.Value) (interface{}, error) {
	message := ""
	url := ""

	// A small test GitLab repo:

	// Testing ways around CORS...

	// gitlab.com
	// WORKS with gitlab.com if browser CORS disabled using a plugin
	message = "https clone of gitlab repo"
	url = "https://gitlab.com/saeedareffard1377666/testproject2.git"

	// github.com
	// WORKS using local-cors-proxy, only for a specific service:
	// yarn global add local-cors-proxy
	// lcp --proxyUrl https://github.com
	url = "http://localhost:8010/proxy/happybeing/p2p-git-portal-poc.git"

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
			repositories[url] = r
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
//	[0] url - remote URL is used to identify a cloned repository (later will be a local identifier)

func listHeadCommits(this js.Value, args []js.Value) (interface{}, error) {
	url := args[0].String()

	println("arg url: ", url)
	url = "http://localhost:8010/proxy/happybeing/p2p-git-portal-poc.git"

	r, found := repositories[url]
	if !found {
		println("Unknown repository: ", url)
		return nil, nil
	}

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
//	[0] url - remote URL is used to identify a cloned repository (later will be a local identifier)
//	[1] first - index of first commit starting at zero
//  [2] last - index of last commit (inclusive, so to return just the first commit first=last=0)

func getHeadCommitsRange(this js.Value, args []js.Value) (interface{}, error) {
	url := args[0].String()
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

	println("arg url: ", url)
	url = "http://localhost:8010/proxy/happybeing/p2p-git-portal-poc.git"

	r, found := repositories[url]
	if !found {
		println("Unknown repository: ", url)
		return nil, nil
	}

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
	commit_index := 0
	total_commits := 0

	err = commitIter.ForEach(func(c *object.Commit) error {
		total_commits += 1
		if commit_index >= first && commit_index <= last {
			commit := make(map[string]interface{}, 0)
			commit["hash"] = c.Hash.String()
			commit["message"] = c.Message

			commits[commit_index] = commit
			commit_index += 1
		}
		return nil
	})

	ret_commits := make(map[string]interface{}, 0)
	ret_commits["length"] = commit_index
	ret_commits["total_commits"] = total_commits
	ret_commits["commits"] = commits

	return ret_commits, nil
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

//// Redundant tests retained temporarily:

func add(this js.Value, args []js.Value) (interface{}, error) {
	ret := 0
	for _, item := range args {
		val := item.Int()
		ret += val
	}

	return ret, nil
}

func err(this js.Value, args []js.Value) (interface{}, error) {
	return nil, errors.New("This is an error")
}

var ready = false

func main() {
	c := make(chan struct{}, 0)

	gobridge.RegisterCallback("uploadFile", uploadFile)
	gobridge.RegisterCallback("listFiles", listFiles)
	gobridge.RegisterCallback("listRepositories", listRepositories)
	gobridge.RegisterCallback("listHeadCommits", listHeadCommits)
	gobridge.RegisterCallback("getHeadCommitsRange", getHeadCommitsRange)
	gobridge.RegisterCallback("testTypes", testTypes)
	gobridge.RegisterCallback("testGitClone", testGitClone)

	gobridge.RegisterCallback("add", add)
	gobridge.RegisterCallback("gitClone", gitClone)
	gobridge.RegisterCallback("raiseError", err)
	gobridge.RegisterValue("someValue", "Hello World")
	gobridge.RegisterValue("numericValue", 123)

	ready = true
	println("Web Assembly is ready")
	<-c // Makes the Go process wait until we want it to end
}
