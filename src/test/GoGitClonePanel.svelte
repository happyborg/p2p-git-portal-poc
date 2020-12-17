<script>
// import { onMount } from 'svelte';
import wasm from '../main.go';

const { cloneRepository } = wasm;

export let disabled = true
export let errorMessage;
export let updateRepositoryUI;

// To use cors-buster proxy locally:
//
// git clone https://github.com/wmhilton/cors-buster
// cd cors-buster && yarn && yarn start
let proxy = "http://localhost:3000/";

let uri;
// Some test URIs
uri = "https://github.com/happybeing/p2p-git-portal-poc";
uri = "https://gitlab.com/saeedareffard1377666/testproject2.git"

async function doGitClone() {
	// console.log("doGitClone() uri, proxy: ", uri, proxy);
	if (uri === undefined || 
		(uri.indexOf('https://') !== 0 && uri.indexOf('http://') !== 0)) {
		errorMessage = "Repository URI must start with http:// or https://"
		return;
	}

	if (proxy !== undefined && proxy !== "" && 
		proxy.indexOf('https://') !== 0 && proxy.indexOf('http://') !== 0) {
		errorMessage = "Proxy URI must start with http:// or https://"
		return;
	}

	let protocol = 'https';
	if (uri.indexOf(protocol + "://") !== 0) {
		protocol = 'http'
	}

	const uriAfterProtocol = uri.substr(protocol.length+3)
	const hostLen = uriAfterProtocol.indexOf('/')
	const host = protocol + "://" + uriAfterProtocol.substr(0,hostLen)
	let path = uriAfterProtocol.substr(hostLen + 1)
	if (path.substr(-1) === '/') path = path.substr(0,path.length-1)

	let proxiedURI = "";
	if (proxy && proxy.indexOf('http') === 0) {
		proxiedURI = proxy;
		if (proxiedURI.substr(-1) !== '/') proxiedURI += '/'
		proxiedURI += uriAfterProtocol
	}

	try {
		await cloneRepository(host, path, proxiedURI, (error) => {
			if (error && error !== "") {
				let message = "clone failed: " + error;
				errorMessage = message;
				console.log(message);
			} else {
				console.log("clone complete");
				updateRepositoryUI(path);
			}
		});
	} catch (e) {
		errorMessage = e
	}
}

</script>
<div>
	<h3>Clone</h3>
	<p>
        <input bind:value={uri} placeholder="URI">
		<button disabled={disabled} type="button" on:click={() => { doGitClone(); }}>Clone</button>
	</p>
	<p>
		<b>Optional:</b><br/>
		CORS proxy: <input bind:value={proxy} placeholder="URI (optional)"><br/>
		Example: http://localhost:3000/<br/>
	</p>
	<p>
		If cloning fails:<br/>
		- try cloning with ".git" suffix (e.g. "https://gitlab.com/anyuser/somerepo.git")<br/>
		- use a CORS proxy or disable CORS with a browser plugin<br/>
	</p>
	<p>
		A proxy works best and you can't clone from github without one.<br/>
	</p>
</div>