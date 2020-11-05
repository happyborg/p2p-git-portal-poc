//+ build js,wasm

package main

import (
	"errors"
	"syscall/js"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/happybeing/webpack-golang-wasm-async-loader/gobridge"
	// "/home/mrh/src/go/src/github.com/happybeing/go-git"
	// "/home/mrh/src/go/src/github.com/happybeing/go-git/storage/memory"
	// "github.com/happybeing/webpack-golang-wasm-async-loader/gobridge"
	// "github.com/happybeing/webpack-golang-wasm-async-loader/gobridge"
	// "/home/mrh/src/wasm/webpack-golang-wasm-async-loader/lib/gobridge"
	// "../../webpack-golang-wasm-async-loader/lib/gobridge"
	// "github.com/go-git/go-git/v5" // with go modules enabled (GO111MODULE=on or outside GOPATH)
	// "github.com/go-git/go-git" // with go modules disabled
)

var global = js.Global()

func gitClone(this js.Value, args []js.Value) (interface{}, error) {
	ret := 0

	println("TODO - implement git clone")
	// Clones the given repository in memory, creating the remote, the local
	// branches and fetching the objects, exactly as:
	println("git clone https://github.com/go-git/go-billy")

	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: "https://github.com/go-git/go-billy",
	})

	// CheckIfError(err)

	// Gets the HEAD history from HEAD, just like this command:
	println("git log")

	// ... retrieves the branch pointed by HEAD
	ref, err := r.Head()
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

func main() {
	c := make(chan struct{}, 0)
	println("Web Assembly is ready")
	gobridge.RegisterCallback("add", add)
	gobridge.RegisterCallback("gitClone", gitClone)
	gobridge.RegisterCallback("raiseError", err)
	gobridge.RegisterValue("someValue", "Hello World")
	gobridge.RegisterValue("numericValue", 123)

	<-c // Makes the Go process wait until we want it to end
}
