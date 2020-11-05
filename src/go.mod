module github.com/happybeing/p2p-git-portal-poc

go 1.13

replace github.com/happybeing/webpack-golang-wasm-async-loader/gobridge => /home/mrh/src/go/src/github.com/happybeing/webpack-golang-wasm-async-loader/gobridge

replace github.com/go-git/go-git/v5 => /home/mrh/src/go/src/github.com/happybeing/go-git

require (
	github.com/go-git/go-git/v5 v5.2.0
	github.com/happybeing/webpack-golang-wasm-async-loader v0.1.0
)
