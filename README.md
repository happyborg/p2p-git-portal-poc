# p2p Git Portal - Proof of Concept

The goal is a github like portal which can be hosted on peer-to-peer or static storage, initially targeting Safe Network. For example, a Safe Network in-browser app using `git-bug` and `git.go` libraries to create a realistic alternative to github.

Using `git-bug` as a library will add support for issues (and later pull requests) to a `git` repo without polluting your commit history, with no need for server side hosting no service provider.

Based on responses on [Mastodon](https://mastodon.technology/@happybeing) there's a lot of interest in a github alternative which is truly p2p, along with decentralisation and open source, which is very encouraging.

**Hashtag [#gitportal](https://mastodon.technology/web/timelines/tag/gitportal)**

## Current Status

### Building a proof-of-concept:
- Svelte UI
- backend with git.go / wasm *in the browser*
- create and access git repo on a BrowserFS file-system

I'll be working on this but help with code or HTML/CSS (see below) is welcome.

After that I'll be trying to make a library version of `git-bug` and increase the functionality to show what's possible, and deploy on peer-to-peer or static web storage.

#### HTML / CSS Work
The app in branch `main` runs and uses default Svelte HTML/CSS styling.

I believe we should change the look and feel to mirror github in the short term. The reason being to help highlight our aim of replicating the key features of a git portal like github.

If you wish to comment on or help with HTML/CSS styling please see https://github.com/happybeing/p2p-git-portal-poc/issues/1.

Later we can differentiate for a more community driven and oriented feature set and our own visual identity, but initially I think it will help to demonstrate the purpose of a proof-of-concept, what is working at any stage and where work is still needed.

At some point we may want to break away from being a github clone because the aims and direction of the two projects are different in key respects, so we should also be thinking about re-inventing the git portal without the goal of centralisation, but of re-sharing the value created by the community in whatever ways the community can benefit.

### Contributions
If you wish to build or contribute to the code, get in touch and I'll add some build instructions as it requires some special setup using my fork of a webpack wasm plugin. Not hard, but not obvious!

## LICENSE

Everything is GPL3.0 unless otherwise stated. Any contributions are accepted on the condition they conform to this license.

See also ./LICENSE