# p2p Git Portal - Proof of Concept

The goal is a github like portal which can be hosted on peer-to-peer or static storage, initially targeting Safe Network. For example, a Safe Network in-browser app using `git-bug` and `git.go` libraries to create a realistic alternative to github.

Using `git-bug` as a library will add support for issues (and later pull requests) to a `git` repo without polluting your commit history, with no need for server side hosting no service provider.

Based on responses on [Mastodon](https://mastodon.technology/@happybeing) there's a lot of interest in a github alternative which is truly p2p, along with decentralisation and open source, which is very encouraging.

**Hashtag [#gitportal](https://mastodon.technology/web/timelines/tag/gitportal)**

## Proof of Concept Summary
### Architecture

- Svelte UI plus Golang/Wasm using Golang wasm webpack plugin
- *in the browser* "backend" with go-git/go-billy and git-bug (Golang/wasm)

### Features
- simple management of repos in browser (on go-billy/memfs)
- ability to upload headless repo into go-billy/memfs
- list repos, visit and interact with a repo
- create and access git repo on a go-billy memfs file-system
- use git-bug in browser to display and edit issues/comments

## Current Status

**WORK IN PROGRESS** - The setup and build instructions below are incomplete so check with me if you want something which builds and can be tested (I have this but you'll need more information).

I'm currently getting the various parts working: go-git and go-billy are available in the browser so I need to tidy this up, prove it is working as needed and add some things to help build on this.

TODO:

[ ] modify HTML/CSS styling. see [#issue 1](https://github.com/happybeing/p2p-git-portal-poc/issues/1).

[ ] fork [happybeing/git-bug](https://github.com/happybeing/git-bug) to make a library version (see git-bug: [issue](https://github.com/happybeing/git-bug/issues/1) and [issue](https://github.com/happybeing/git-bug/issues/2))

[ ] add happybeing/git-bug to the poc compilation (as with happybeing/go-git) (see [issue #2](https://github.com/happybeing/p2p-git-portal-poc/issues/2))

[ ] ACTIVE WIP: add file upload UI to the browser, storing files using go-billy/memfs

[ ] add instructions for setting up local forks of go-git, go-billy and git-bug

[ ] specify/build UI to test go-git repo operations on git repo uploaded to go-billy/memfs

#### HTML / CSS Work
The app in branch `main` runs and uses default Svelte HTML/CSS styling.

I believe we should change the look and feel to mirror github in the short term. The reason being to help highlight our aim of replicating the key features of a git portal like github.

If you wish to comment on or help with HTML/CSS styling please see [#issue 1](https://github.com/happybeing/p2p-git-portal-poc/issues/1).

Later we can differentiate for a more community driven and oriented feature set and our own visual identity, but initially I think it will help to demonstrate the purpose of a proof-of-concept, what is working at any stage and where work is still needed.

At some point we may want to break away from being a github clone because the aims and direction of the two projects are different in key respects, so we should also be thinking about re-inventing the git portal without the goal of centralisation, but of re-sharing the value created by the community in whatever ways the community can benefit.

## Development Setup

### Pre-requisites
Tested using Golang v1.13, `node` v14.14 and `yarn` v1.22. You could use `npm` for linking but I use yarn for this as it seems more reliable.

### Clone the plugin and this template

**NOTE** - these instructions are incomplete - get in touch if you want something which works.

TODO: add instructions for setting up local forks of go-git, go-billy and git-bug.

Note: the directories need to be adjacent in order to work as is. If not you will have to modify the path to `gobridge` in your `main.go` import. I'm not sure why this is needed but maybe the plugin wasn't finalised.

If you have `node` and `yarn`, on Linux you should be able to just copy the following and paste it into your terminal.
``` bash
# Make sure GOROOT and GOPATH are set in the terminal, for example:
export GOROOT=`go env GOROOT`
export GOPATH=`go env GOPATH`

# Get the plugin
mkdir -p ~/src/wasm_test
cd ~/src/wasm_test
git clone https://github.com/happybeing/webpack-golang-wasm-async-loader
cd webpack-golang-wasm-async-loader
npm install && npm run build
yarn link

# Get the app template
cd ~/src/wasm_test
git clone https://github.com/happybeing/p2p-git-portal-poc
cd  p2p-git-portal-poc
yarn link golang-wasm-async-loader
yarn && yarn build

# If all looks good, start the app
yarn dev
# Open app in the browser by visiting localhost:5000
```
The plugin directory and its package name are different, so you need to use the `yarn link` command as shown.

### Contributions
If you wish to build or contribute to the code, get in touch and I'll add some build instructions as it requires some special setup using my fork of a webpack wasm plugin. Not hard, but not obvious!

## LICENSE

Everything is GPL3.0 unless otherwise stated. Any contributions are accepted on the condition they conform to this license.

See also ./LICENSE