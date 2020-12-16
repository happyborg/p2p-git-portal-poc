//+ build js,wasm

package main

import (
	"fmt"
	"log"
	"syscall/js"

	"github.com/happybeing/p2p-git-portal-poc/src/repo"
	"github.com/happybeing/webpack-golang-wasm-async-loader/gobridge"
)

var global = js.Global()

func uploadFile(this js.Value, args []js.Value) (interface{}, error) {
	// if !ready {
	// 	return nil, nil
	// }

	ret := 0
	fullPath := args[0].String()

	array := args[1]
	// println("Array byteLength: ", array.Get("byteLength").Int())
	buf := make([]byte, array.Get("byteLength").Int())
	n := js.CopyBytesToGo(buf, array)

	fmt.Println("GO uploading: ", fullPath, n)

	dst, err := repo.Filesystem.Create(fullPath)
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

// args[]:
//	[0] path - local repo path to files in a repository
func getDirectoryEntries(this js.Value, args []js.Value) (interface{}, error) {
	path := args[0].String()
	return listFiles(path)
}

func listFiles(path string) ([]interface{}, error) {
	println("listFiles()")
	listing, err := repo.Filesystem.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	entries := make([]interface{}, len(listing))
	for index, f := range listing {
		entry := make(map[string]interface{}, 0)
		entry["name"] = f.Name()
		entry["modified"] = f.ModTime().UTC().Unix()
		entry["mode"] = f.Mode().String() // FileMode e.g. "drwxrwxrwx"
		entry["size"] = f.Size()          // int64
		entry["sys"] = f.Sys()            // underlying data source (can return nil)
		if f.IsDir() {
			entry["type"] = "directory"
		} else {
			entry["type"] = "file"
		}
		entries[index] = entry
	}

	return entries, err
}

//// git-bug gogit.Repository tests

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
	gobridge.RegisterCallback("listHeadCommits", repo.ListHeadCommits)
	gobridge.RegisterCallback("testTypes", testTypes)
	gobridge.RegisterCallback("testGitClone", repo.GitCloneTest)

	gobridge.RegisterCallback("openRepository", repo.OpenRepository)
	gobridge.RegisterCallback("cloneRepository", repo.CloneRepository)
	gobridge.RegisterCallback("getRepositoryList", repo.GetRepositoryList)
	gobridge.RegisterCallback("getHeadCommitsRange", repo.GetHeadCommitsRange)
	gobridge.RegisterCallback("getIssuesForRepo", repo.GetIssuesForRepo)
	gobridge.RegisterCallback("getDirectory", getDirectoryEntries)

	gobridge.RegisterCallback("newRepository", repo.NewRepository)
	ready = true
	println("Web Assembly is ready")
	<-c // Makes the Go process wait until we want it to end
}
