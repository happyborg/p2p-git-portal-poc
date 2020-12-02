# p2p Git Portal - Proof of Concept

The goal is a github like portal which can be hosted on peer-to-peer or static storage, initially targeting Safe Network. For example, a Safe Network in-browser app using `git-bug` and `go-git` libraries to create a realistic alternative to github.

Using `git-bug` as a library will add support for issues (and later pull requests) to a `git` repo without polluting your commit history, with no need for server side hosting no service provider.

Based on responses on [Mastodon](https://mastodon.technology/@happybeing) there's a lot of interest in a github alternative which is truly p2p, along with decentralisation and open source, which is very encouraging.

**Hashtag #gitportal:** [Mastodon](https://mastodon.technology/web/timelines/tag/gitportal) | [Twitter](https://twitter.com/hashtag/gitportal)

**Chat (matrix/element):** [#gitportal](https://riot.im/app/#/room/#gitportal:matrix.org)

**Forum (discourse):** [Git Portal Discussion](https://safenetforum.org/t/safenetwork-git-portal-discussion/32793?u=happybeing)

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

**WORK IN PROGRESS** - everything in the above list is working apart for the last item: use git-bug to display and edit issues/comments.

Remaining activity for the proof-of-concept:

TODO:

[ ] modify HTML/CSS styling. see [#issue 1](https://github.com/happybeing/p2p-git-portal-poc/issues/1).

[ ] fork [happybeing/git-bug](https://github.com/happybeing/git-bug) to make a library version (see git-bug: [issue](https://github.com/happybeing/git-bug/issues/1) and [issue](https://github.com/happybeing/git-bug/issues/2))

[ ] add happybeing/git-bug to the poc compilation (as with happybeing/go-git) (see [issue #2](https://github.com/happybeing/p2p-git-portal-poc/issues/2))

[ ] add instructions for setting up local forks of go-git, go-billy and git-bug

#### HTML / CSS Work
The app in branch `main` runs and uses default Svelte HTML/CSS styling.

I believe we should change the look and feel to mirror github in the short term. The reason being to help highlight our aim of replicating the key features of a git portal like github.

If you wish to comment on or help with HTML/CSS styling please see [#issue 1](https://github.com/happybeing/p2p-git-portal-poc/issues/1).

Later we can differentiate for a more community driven and oriented feature set and our own visual identity, but initially I think it will help to demonstrate the purpose of a proof-of-concept, what is working at any stage and where work is still needed.

At some point we may want to break away from being a github clone because the aims and direction of the two projects are different in key respects, so we should also be thinking about re-inventing the git portal without the goal of centralisation, but of re-sharing the value created by the community in whatever ways the community can benefit.

## Setup

### Pre-requisites

Install Golang v1.15, `node` v14.14 and `yarn` v1.22. You could use `npm`, I just prefer `yarn`.

I recommend using `nvm` (node version manager to install `node`) and `gvm` (Go version manager to install Golang) as it makes it easier to install, upgrade and switch between versions of these dependencies.

### Get the Code

If you have `node` and `yarn` installed, on Linux you should be able to just copy the following and paste it into your terminal.
``` bash
# Make sure GOROOT and GOPATH are set in the terminal, for example:
export GOROOT=`go env GOROOT`
export GOPATH=`go env GOPATH`

# Get the app
mkdir -p ~/src/go_wasm
cd ~/src/go_wasm
git clone https://github.com/happybeing/p2p-git-portal-poc
cd  p2p-git-portal-poc
yarn && yarn build

# If all looks good, start the app
yarn dev
# Open app in the browser by visiting localhost:8080
```

### Setup a CORS Proxy
To clone repositories you will need a way to overcome CORS errors because we're working entirely within the web browser. You can try disabling this with a browser plugin, and this works with some services but not github.

So I recommend using a CORS proxy. You can set one up locally as follows:
```bash
	git clone https://github.com/wmhilton/cors-buster
	cd cors-buster && yarn && yarn start
```

### Testing
The above is all that's needed to test and for most development.

Make sure you have the proxy running and make sure the proxy URI is set correctly in the UI. Now click "Clone". 

A sample repository is already set in the UI, and should be cloned in a few seconds. Large repositories will of course take much longer, so be prepared to wait!

When cloning is finished this the list of repositories on the top left, and the list of commits on the right will be updated. If you have more than one repository cloned, you can click on the list (top left) to show the commits for that repository.

If you have problems, open the web browser console to look for any error messages and feel free to open an issue if you have difficulties.

Note: drag and drop of files is not properly supported yet.

### Screenshot
<img src="./gitportal-poc-stage1.png"/>

### Development
Hot reloading generally works well, but if you have problems getting rid of a compilation error after you think you've fixed it, restart the `yarn dev` command.

If you want to make changes to the go-git/go-billy, git-bug or the Go/wasm plugin, you need to be set up with local versions as described next, but you should not need to do this.

### Code Documentation
You can view **developer documentation** for local code and package dependencies in a web browser, showing functions, types etc. for each package.

To view developer documentation:

- Install `godoc` using `go get` / `go install`. Linux users may install directly from their package repositories but this may not be as up-to-date (e.g. `sudo apt install godoc`).

In the directory next to the `go.mod` of the code you wish to browser (e.g. `./p2p-git-portal-poc/src`), run the following but be patient as first time it takes a while to generate the docs and doesn't show any progress in the terminal.
```bash
# This updates the documentation and then acts as a web server
godoc -http=":6060"
```
Then visit `http://localhost:6060` and click 'Packages' and search the page for the package you want to browse (e.g. `git-bug`)


### Contributions
If you wish to build or contribute to the code, get in touch and I'll add some build instructions as it requires some special setup using my fork of a webpack wasm plugin. Not hard, but not obvious!

## LICENSE

Everything is GPL3.0 unless otherwise stated. Any contributions are accepted on the condition they conform to this license.

See also ./LICENSE