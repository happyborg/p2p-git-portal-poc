//+ build js,wasm

package main

import (
	"errors"
	"fmt"
	"log"
	"syscall/js"

	"github.com/happybeing/webpack-golang-wasm-async-loader/gobridge"
	// "/home/mrh/src/go/src/github.com/happybeing/webpack-golang-wasm-async-loader/gobridge"

	// OK FOR CLI if I have gobridge/go.mod containing:
	// module github.com/happybeing/webpack-golang-wasm-async-loader/gobridge
	// go 1.13

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	//
	// TRY FOR yarn test:
	// "github.com/happybeing/go-git"
	// "github.com/happybeing/go-git/storage/memory"
	// Gets deleted on save"github.com/happybeing/go-billy"
	//
	// TRY FOR yarn build:
	// "github.com/happybeing/go-git/v5"
	// "github.com/happybeing/go-git/v5/storage/memory"
	// //
	// "github.com/happybeing/webpack-golang-wasm-async-loader/gobridge"
	// "github.com/happybeing/webpack-golang-wasm-async-loader/gobridge"
	// "/home/mrh/src/wasm/webpack-golang-wasm-async-loader/lib/gobridge"
	// "../../webpack-golang-wasm-async-loader/lib/gobridge"
	// "github.com/go-git/go-git/v5" // with go modules enabled (GO111MODULE=on or outside GOPATH)
	// "github.com/go-git/go-git" // with go modules disabled
)

var global = js.Global()

var fs = memfs.New()

func upload(this js.Value, files []js.Value) (interface{}, error) {
	// if !ready {
	// 	return nil, nil
	// }

	ret := 0

	for _, file := range files {
		println("GO uploading: ", file.String())
		// if file.IsDir() {
		// 	continue
		// }

		// src, err := origin.Open(file)
		// if err != nil {
		// 	return nil, err
		// }

		dst, err := fs.Create(file.String())
		if err != nil {
			return nil, err
		}

		// if _, err = io.Copy(dst, src); err != nil {
		// 	return nil, err
		// }

		if err := dst.Close(); err != nil {
			return nil, err
		}

		// if err := src.Close(); err != nil {
		// 	return nil, err
		// }
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

func gitClone(this js.Value, args []js.Value) (interface{}, error) {
	ret := 0

	url := ""
	message := ""
	println("TODO - implement git clone")
	// Clones the given repository in memory, creating the remote, the local
	// branches and fetching the objects, exactly as:
	// println("git clone https://github.com/happybeing/p2p-git-portal-poc")
	// r, _ := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
	// 	URL: "https://github.com/happybeing/p2p-git-portal-poc",
	// })

	// println("git://github.com/happybeing/p2p-git-portal-poc")
	// r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
	// 	URL: "git://github.com/happybeing/p2p-git-portal-poc",
	// })

	message = "ssh server git.Clone() "
	url = "mrh@127.0.0.1:.gitserver/test-git"
	// url = "mrh@127.0.0.1:home/mrh/.gitserver/test-git"

	// message = "access token auth / git clone "
	// url = "https://github.com/happybeing/p2p-git-portal-poc"
	// url = "git@github.com:happybeing/p2p-git-portal-poc"
	// // url = "git://github.com/happybeing/p2p-git-portal-poc"
	// user := "happybeing"
	// token := "248a93a240e309a9246c1caed2cd094bb1c09e70"

	// message = "local server git.Clone() "
	// url = "http://127.0.0.1:1236"

	println(message, url)

	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		// The intended use of a GitHub personal access token is in replace of your password
		// because access tokens can easily be revoked.
		// https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
		// Auth: &http.BasicAuth{
		// 	Username: user, // yes, this can be anything except an empty string
		// 	Password: token,
		// },
		URL: url,
	})

	if err != nil {
		println("git.Clone() failed: ", err.Error())
	}
	// CheckIfError(err)

	// Gets the HEAD history from HEAD, just like this command:
	println("git log")

	// ... retrieves the branch pointed by HEAD
	_, err2 := r.Head()
	if err2 != nil {
		println("r.Head() failed: ", err2.Error())
	}
	// ref, err := r.Head()
	// CheckIfError(err)

	return ret, nil
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

	gobridge.RegisterCallback("upload", upload)
	gobridge.RegisterCallback("listFiles", listFiles)

	gobridge.RegisterCallback("add", add)
	gobridge.RegisterCallback("gitClone", gitClone)
	gobridge.RegisterCallback("raiseError", err)
	gobridge.RegisterValue("someValue", "Hello World")
	gobridge.RegisterValue("numericValue", 123)

	ready = true
	println("Web Assembly is ready")
	<-c // Makes the Go process wait until we want it to end
}
