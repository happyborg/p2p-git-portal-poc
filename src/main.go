//+ build js,wasm

package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
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

	// message = "ssh server git.Clone() "
	// url = "mrh@127.0.0.1:.gitserver/test-git"
	// url = "mrh@127.0.0.1:home/mrh/.gitserver/test-git"

	// message = "access token auth / git clone "
	// url = "https://github.com/happybeing/p2p-git-portal-poc"
	// url = "git@github.com:happybeing/p2p-git-portal-poc"
	// // url = "git://github.com/happybeing/p2p-git-portal-poc"
	// user := "happybeing"
	// token := oops ""

	// message = "expriment with github URIs"
	// url = "https://github.com/happybeing/p2p-git-portal-poc.git"

	// message = "local server git.Clone() "
	// url = "http://127.0.0.1:1236"

	// message = "file:// URI git.Clone() "
	// url = "file:///home/mrh/.gitserver/testrepo"

	println(message, url)

	url = "https://gitlab.com/weblate/libvirt"
	fs := memfs.New()
	storage := memory.NewStorage()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		_, err := git.Clone(storage, fs, &git.CloneOptions{
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

		wg.Done()
	}()
	wg.Wait()

	// CheckIfError(err)

	// Gets the HEAD history from HEAD, just like this command:
	println("git log")

	// ... retrieves the branch pointed by HEAD
	// _, err2 := r.Head()
	// if err2 != nil {
	// 	println("r.Head() failed: ", err2.Error())
	// }
	// ref, err := r.Head()
	// CheckIfError(err)

	return ret, nil
}

/*
// Disabling CORS (or using proxy) with GitLab
// message = "https clone of gitlab repos"
// url = "https://gitlab.com/weblate/libvirt"
almost works:
- does OPTIONS then GET https://gitlab.com/weblate/libvirt/info/refs?service=git-upload-pack
- is redirected to  https://gitlab.com/weblate/libvirt.git/info/refs?service=git-upload-pack
- does OPTIONS then GET  https://gitlab.com/weblate/libvirt.git/info/refs?service=git-upload-pack
- receives payload of 66KB
- browser console error:
	https clone of gitlab repos https://gitlab.com/weblate/libvirt wasm_exec.js:45
	fatal error: all goroutines are asleep - deadlock! wasm_exec.js:45
	<empty string> wasm_exec.js:45
	goroutine 1 [chan receive]: wasm_exec.js:45
	main.main() wasm_exec.js:45
	/home/mrh/src/wasm/p2p-git-portal-poc/src/main.go:197 +0x12
[ ] looks to me like the requests are issued and Go crashes BEFORE the response
	because the console errors are all printed and *then* the HTTP output
	follows over a second or two.
	-> Looks like [issue 41310](https://github.com/golang/go/issues/41310)
	"You have to create a separate event handler loop or some other mechanism
	to make http requests and read responses."
	[ ] asked how to solve this with go-git [here](https://github.com/golang/go/issues/41310#issuecomment-725475075)
		[x] try wrapping repo.Clone() call in a separate go routine
			-> gets further (looks like the Clone() happens) but same error at end
	[ ] will not solve issue with github which fails on the requests with CORS errors even with CORS disabled in browser!
*/
func testGitCloneDeadlock(this js.Value, args []js.Value) (interface{}, error) {

	// Disabling CORS (or using proxy) with GitLab
	message := "https clone of gitlab repos"
	url := "https://gitlab.com/weblate/libvirt"

	println(message, url)

	url = "https://gitlab.com/weblate/libvirt"
	fs := memfs.New()
	storage := memory.NewStorage()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		_, err := git.Clone(storage, fs, &git.CloneOptions{URL: url})
		if err != nil {
			println("git.Clone() failed: ", err.Error())
		}

		wg.Done()
	}()
	wg.Wait()

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
	gobridge.RegisterCallback("testGitCloneDeadlock", testGitCloneDeadlock)

	gobridge.RegisterCallback("add", add)
	gobridge.RegisterCallback("gitClone", gitClone)
	gobridge.RegisterCallback("raiseError", err)
	gobridge.RegisterValue("someValue", "Hello World")
	gobridge.RegisterValue("numericValue", 123)

	ready = true
	println("Web Assembly is ready")
	<-c // Makes the Go process wait until we want it to end
}
