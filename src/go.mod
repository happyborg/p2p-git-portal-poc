module github.com/happybeing/p2p-git-portal-poc

go 1.13

require (
	github.com/MichaelMure/git-bug v0.0.0-20201123161810-b9961535505b
	github.com/go-git/go-billy/v5 v5.0.0
	github.com/go-git/go-git v4.7.0+incompatible // indirect
	github.com/go-git/go-git/v5 v5.2.0
	github.com/happybeing/webpack-golang-wasm-async-loader/gobridge v0.0.0-20201125123737-917b0db1c9a8
// github.com/happybeing/webpack-golang-wasm-async-loader/gobridge v0.0.0-20201121121723-6f1ca5446694
// gopkg.in/src-d/go-git.v4 v4.13.1 // indirect
)

// Use forks with Go/wasm support
replace github.com/go-git/go-git/v5 => github.com/happybeing/go-git/v5 v5.2.1

replace github.com/go-git/go-billy/v5 => github.com/happybeing/go-billy/v5 v5.0.1

// For development with local go-git/go-billy and wasm plugin:
replace github.com/MichaelMure/git-bug => /home/mrh/src/go/src/github.com/happybeing/git-bug

replace github.com/99designs/keyring => /home/mrh/src/go/src/github.com/happybeing/keyring

// Dependencies of git-bug which require tweaks for wasm
// replace github.com/99designs/keyring => /home/mrh/src/go/src/github.com/happybeing/keyring

// replace github.com/go-git/go-git/v5 => /home/mrh/src/go/src/github.com/happybeing/go-git
// replace github.com/go-git/go-billy/v5 => /home/mrh/src/go/src/github.com/happybeing/go-billy

// To develop with a local version of the plugin you need to:
// - uncomment the following 'replace'
// - `yarn link` in your local development copy of the plugin
// - `yarn link golang-wasm-async-loader2` in ./p2p-git-portal-poc
// replace github.com/happybeing/webpack-golang-wasm-async-loader/gobridge => /home/mrh/src/go/src/github.com/happybeing/webpack-golang-wasm-async-loader/gobridge
