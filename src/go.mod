module github.com/happybeing/p2p-git-portal-poc

go 1.13

require (
	github.com/go-git/go-git v4.7.0+incompatible // indirect
	// Try for yarn test:
	// github.com/happybeing/go-git/v5
	// github.com/happybeing/go-billy/v5
	// OK for CLI:
	github.com/go-git/go-git/v5 v5.2.0
	github.com/happybeing/webpack-golang-wasm-async-loader/gobridge v0.1.0
	gopkg.in/src-d/go-git.v4 v4.13.1 // indirect
)

// OK for CLI:
replace github.com/go-git/go-git/v5 => /home/mrh/src/go/src/github.com/happybeing/go-git

replace github.com/go-git/go-billy/v5 => /home/mrh/src/go/src/github.com/happybeing/go-billy

replace github.com/happybeing/webpack-golang-wasm-async-loader/gobridge => /home/mrh/src/go/src/github.com/happybeing/webpack-golang-wasm-async-loader/gobridge

// Try for yarn test:

// replace github.com/go-git/go-git => /home/mrh/src/go/src/github.com/happybeing/go-git

// replace github.com/go-git/go-billy => /home/mrh/src/go/src/github.com/happybeing/go-billy

// replace github.com/happybeing/webpack-golang-wasm-async-loader/gobridge => /home/mrh/src/go/src/github.com/happybeing/webpack-golang-wasm-async-loader/gobridge
