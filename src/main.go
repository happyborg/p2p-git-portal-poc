//+ build js,wasm

package main

import (
	"fmt"
	"log"
	"syscall/js"

	"github.com/happybeing/p2p-git-portal-poc/src/repo"
	"github.com/happybeing/webpack-golang-wasm-async-loader/gobridge"

	// "github.com/happybeing/p2p-git-portal-poc/src/repo"
	// repo "./repo"

	// OK FOR CLI if I have gobridge/go.mod containing:
	// module github.com/happybeing/webpack-golang-wasm-async-loader/gobridge
	// go 1.13

	"github.com/go-git/go-billy/v5/memfs"

	"github.com/MichaelMure/git-bug/cache"
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

//// git-bug gogit.Repository tests

func testGitBug(this js.Value, args []js.Value) (interface{}, error) {
	cache := cache.NewMultiRepoCache()
	if cache == nil {
		println("testGitBug() FAILED to create cache")
	} else {
		println("testGitBug() created multi-repo cache!", cache)
	}
	return nil, nil
}

func testRepoInit(this js.Value, args []js.Value) (interface{}, error) {
	println("testRepoInit()...")
	return nil, repo.PocRepoInitialise()
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

////// Go/wasm initialisation

var ready = false

func main() {
	c := make(chan struct{}, 0)

	gobridge.RegisterCallback("uploadFile", uploadFile)
	gobridge.RegisterCallback("listFiles", listFiles)
	gobridge.RegisterCallback("listHeadCommits", repo.ListHeadCommits)
	gobridge.RegisterCallback("testTypes", testTypes)
	gobridge.RegisterCallback("testGitClone", repo.GitCloneTest)

	gobridge.RegisterCallback("cloneRepository", repo.CloneRepository)
	gobridge.RegisterCallback("getRepositoryList", repo.GetRepositoryList)
	gobridge.RegisterCallback("getHeadCommitsRange", repo.GetHeadCommitsRange)

	// gobridge.RegisterCallback("testGitBug", testGitBug)
	gobridge.RegisterCallback("testGitBug", testRepoInit)
	ready = true
	println("Web Assembly is ready")
	<-c // Makes the Go process wait until we want it to end
}
