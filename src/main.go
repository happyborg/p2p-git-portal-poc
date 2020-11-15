//+ build js,wasm

package main

import (
	"errors"
	"fmt"
	"log"
	"syscall/js"

	"github.com/happybeing/webpack-golang-wasm-async-loader/gobridge"

	// OK FOR CLI if I have gobridge/go.mod containing:
	// module github.com/happybeing/webpack-golang-wasm-async-loader/gobridge
	// go 1.13
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
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
