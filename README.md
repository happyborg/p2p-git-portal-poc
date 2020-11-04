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

After that I'll be trying to make a library version of `git-bug` and increase the functionality to show what's possible, and deploy on peer-to-peer or static web storage.

### Contributions
If you wish to contribute to the code, get in touch and I'll add some build instructions as it requires some special setup using my fork of a webpack wasm plugin. Not hard, but not obvious!

## LICENSE

Everything is GPL3.0 unless otherwise stated. Any contributions are accepted on the condition they conform to this license.

See also ./LICENSE