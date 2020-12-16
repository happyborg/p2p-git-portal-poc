# p2p Git Portal - Proof of Concept

Git Portal is a github style sharing and collaboration application without server side code (no-backend).

- The aim is to provide an alternative to the centralisation of github and similar third party hosted services, and to avoid the need to self host such a service.
- The approach is to provide web front end in a static application which can run in a browser directly from static or peer-to-peer storage without the need for server side code.
- The benefit of this will be the convenience of a web application, without dependence on a third party and avoiding the vulnerabilities associated with servers and centralised services.

Based on responses on [Mastodon](https://mastodon.technology/@happybeing) there's a lot of interest in a github alternative which is truly p2p, along with decentralisation and open source, which is very encouraging.

**Hashtag #gitportal:** [Mastodon](https://mastodon.technology/web/timelines/tag/gitportal) | [Twitter](https://twitter.com/hashtag/gitportal)

**Chat (matrix/element):** [#gitportal](https://riot.im/app/#/room/#gitportal:matrix.org)

**Forum (discourse):** [Git Portal Discussion](https://safenetforum.org/t/safenetwork-git-portal-discussion/32793?u=happybeing)

## Screenshot
<img src="./gitportal-poc-screenshot.png"/>

## Objective
The initial objective establish feasibility this project will create an initial app designed to run from static storage or on the peer-to-peer [Safe Network](https://safenetwork.tech) (a secure, anonymous, private peer-to-peer storage and communications platform). The result could then be adapted for use on any decentralised storage, and be the basis for more ambitious projects.
## Approach

The implementation identified `git-bug` (which adds issues to git functionality without polluting the commit history) and `go-git`/`go-billy` on which it depends. These are ideal for the proof-of-concept because they offer a good solution to providing issues on top of git, and a rapid route to prove the concept with a small team.

Later the implementation will be reviewed in order to improve the technical solution. 

Once we have a proof-of-concept application, a larger vision can be co-created which will align with the values and wants of collaborative software developers, as opposed to those of the Microsoft which owns and runs github for profit.
### Golang v Rust

One technical change to consider is moving from Golang/Wasm to Rust/Wasm to improve security and performance. Golang was chosen mainly for convenience and speed, but is slow, has a large runtime (12MB), and is not the best choice for building secure applications. All things which Rust/Wasm does very well.

## Proof of Concept Summary
### Architecture

- Svelte UI plus Golang/Wasm using Golang wasm webpack plugin
- *in the browser* 'no-backend' with go-git/go-billy and git-bug (Golang/wasm)

### Features
- initialise, clone or upload repository 
- list repositories
- select and interact with a repository
- browse the worktree, commits and issues
- create/view/edit issues and comments

## Current Status

Most of the above feautures are either working or proven to be feasible by the work done so far. The UI is still a test harness but work has begun (mid December 2020) to provide a more github like look and feel.

Remaining activity for the proof-of-concept:

### Activity & Opportunities to Help

- **Modify Styling** of HTML/CSS. **Status:** plenty to do!

	**Help Wanted** with tasks suitable for FOSS folks with anything from basic CSS skills and a desire to learn some Svelte (with my help) to folks who want to create the vision for a future git portal UX (user experience). Some notes are included below and there's a discussion on [#issue 1](https://github.com/happybeing/p2p-git-portal-poc/issues/1).

- **Improve Features and UI** to create a more realistic demonstration and formulate a suitable feature set for the API to support the app. 

	**Help Wanted** with several small tasks in Svelte front-end and Golang no-server back-end, but these are not yet written down. So ask if you may be able to help and I will begin to turn my thoughts into tasks with your prompting. Bigger tasks will flow from this if you want to help design or implement the features and API providing these to the front-end.

- **Modify git-bug** to work with a billy.Filesystem and compile to web assembly. **Status:** I got this! The fork [happybeing/git-bug](https://github.com/happybeing/git-bug) is sufficient for remaining work to complete this proof-of-concept. [@happybeing](https://github.com/happybeing) will bring it up to date with the latest git-bug and work with git-bug author [@MichaelMure](https://github.com/MichaelMure) to provide a proper API for use by apps, as summarised in this [issue comment](https://github.com/happybeing/git-bug/issues/2#issuecomment-742494498).

#### HTML / CSS Work
(Begun mid Deceber 2020)

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

### Development
Hot reloading generally works well, but if you have problems getting rid of a compilation error after you think you've fixed it, restart the `yarn dev` command.

If you want to make changes to the Go/wasm plugin or to any go packages used by this project, you must set up local versions under your `go/src` directory and modify your `./go.mod` to use these rather than download them. For example, my `./go.mod` includes 'replace' statements for each package I'm working on locally:

```golang
replace github.com/MichaelMure/git-bug => /home/mrh/src/go/src/github.com/happybeing/git-bug

```

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

See also [./LICENSE](./LICENSE)