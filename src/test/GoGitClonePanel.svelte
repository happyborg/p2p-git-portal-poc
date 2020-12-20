<script>
// import { onMount } from 'svelte';
import wasm from '../main.go';

const { cloneRepository } = wasm;

let visible = false;
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
<hr>
<div class="my-3">
	<label class="inline-block text-gray-800 text-m font-semibold mb-2 mt-8" for="">Clone repository</label> 
	<div class="relative">
		<button on:click={() => visible = !visible} class="h-5 w-5 rounded-full overflow-hidden focus:outline-none absolute right-0 -top-7">
		<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M7.429 1.525a6.593 6.593 0 011.142 0c.036.003.108.036.137.146l.289 1.105c.147.56.55.967.997 1.189.174.086.341.183.501.29.417.278.97.423 1.53.27l1.102-.303c.11-.03.175.016.195.046.219.31.41.641.573.989.014.031.022.11-.059.19l-.815.806c-.411.406-.562.957-.53 1.456a4.588 4.588 0 010 .582c-.032.499.119 1.05.53 1.456l.815.806c.08.08.073.159.059.19a6.494 6.494 0 01-.573.99c-.02.029-.086.074-.195.045l-1.103-.303c-.559-.153-1.112-.008-1.529.27-.16.107-.327.204-.5.29-.449.222-.851.628-.998 1.189l-.289 1.105c-.029.11-.101.143-.137.146a6.613 6.613 0 01-1.142 0c-.036-.003-.108-.037-.137-.146l-.289-1.105c-.147-.56-.55-.967-.997-1.189a4.502 4.502 0 01-.501-.29c-.417-.278-.97-.423-1.53-.27l-1.102.303c-.11.03-.175-.016-.195-.046a6.492 6.492 0 01-.573-.989c-.014-.031-.022-.11.059-.19l.815-.806c.411-.406.562-.957.53-1.456a4.587 4.587 0 010-.582c.032-.499-.119-1.05-.53-1.456l-.815-.806c-.08-.08-.073-.159-.059-.19a6.44 6.44 0 01.573-.99c.02-.029.086-.075.195-.045l1.103.303c.559.153 1.112.008 1.529-.27.16-.107.327-.204.5-.29.449-.222.851-.628.998-1.189l.289-1.105c.029-.11.101-.143.137-.146zM8 0c-.236 0-.47.01-.701.03-.743.065-1.29.615-1.458 1.261l-.29 1.106c-.017.066-.078.158-.211.224a5.994 5.994 0 00-.668.386c-.123.082-.233.09-.3.071L3.27 2.776c-.644-.177-1.392.02-1.82.63a7.977 7.977 0 00-.704 1.217c-.315.675-.111 1.422.363 1.891l.815.806c.05.048.098.147.088.294a6.084 6.084 0 000 .772c.01.147-.038.246-.088.294l-.815.806c-.474.469-.678 1.216-.363 1.891.2.428.436.835.704 1.218.428.609 1.176.806 1.82.63l1.103-.303c.066-.019.176-.011.299.071.213.143.436.272.668.386.133.066.194.158.212.224l.289 1.106c.169.646.715 1.196 1.458 1.26a8.094 8.094 0 001.402 0c.743-.064 1.29-.614 1.458-1.26l.29-1.106c.017-.066.078-.158.211-.224a5.98 5.98 0 00.668-.386c.123-.082.233-.09.3-.071l1.102.302c.644.177 1.392-.02 1.82-.63.268-.382.505-.789.704-1.217.315-.675.111-1.422-.364-1.891l-.814-.806c-.05-.048-.098-.147-.088-.294a6.1 6.1 0 000-.772c-.01-.147.039-.246.088-.294l.814-.806c.475-.469.679-1.216.364-1.891a7.992 7.992 0 00-.704-1.218c-.428-.609-1.176-.806-1.82-.63l-1.103.303c-.066.019-.176.011-.299-.071a5.991 5.991 0 00-.668-.386c-.133-.066-.194-.158-.212-.224L10.16 1.29C9.99.645 9.444.095 8.701.031A8.094 8.094 0 008 0zm1.5 8a1.5 1.5 0 11-3 0 1.5 1.5 0 013 0zM11 8a3 3 0 11-6 0 3 3 0 016 0z"></path></svg></button>
		{#if visible}
		  <div class="absolute right-0 w-300 mt-1 p-4 bg-white border rounded shadow-xl z-30">   
			  <div class="text-sm">
				  <p>
					  <b>CORS proxy:</b><br/>
					<input class="appearance-none w-full p-3 mb-2 text-xs font-semibold leading-none bg-gray-100 rounded outline-none" bind:value={proxy} placeholder="URI (optional)"><br/>
					  Example: http://localhost:3000/<br/>
				  </p>
				  <div class="mt-2">
					  <strong>If cloning fails:</strong>
					  <ul class="list-disc list-ouside ml-4">
						  <li>try cloning with ".git" suffix (e.g. "https://gitlab.com/anyuser/somerepo.git")</li>
						  <li>use a CORS proxy or disable CORS with a browser plugin</li>
					  </ul>
					</div>
				  <p class="mt-2 text-xs text-left text-gray-600 bg-blue-50 border border-blue-100 p-2 rounded-md">
					  A proxy works best and you can't clone from github without one.<br/>
				  </p>
			  </div>
		  </div>
	  {/if}
	  </div>
	  <!-- // Dropdown -->
	<input bind:value={uri} class="appearance-none w-full p-3 text-xs font-semibold leading-none bg-gray-100 rounded outline-none" type="text" name="field-name" placeholder="Repository url">

	<p class="mt-3">
		<button on:click={() => { doGitClone(); }}  class="bg-white text-sm text-gray-800 font-bold border border-gray-200 hover:border-red-600 hover:bg-red-500 hover:text-white py-2 px-5 inline-flex items-center">
			<span class="mr-2">Clone</span>
			<svg xmlns="http://www.w3.org/2000/svg" fill="currentColor" viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M15 0H9v7c0 .55.45 1 1 1h1v1h1V8h3c.55 0 1-.45 1-1V1c0-.55-.45-1-1-1zm-4 7h-1V6h1v1zm4 0h-3V6h3v1zm0-2h-4V1h4v4zM4 5H3V4h1v1zm0-2H3V2h1v1zM2 1h6V0H1C.45 0 0 .45 0 1v12c0 .55.45 1 1 1h2v2l1.5-1.5L6 16v-2h5c.55 0 1-.45 1-1v-3H2V1zm9 10v2H6v-1H3v1H1v-2h10zM3 8h1v1H3V8zm1-1H3V6h1v1z"></path></svg>
		  </button>
	</p>
</div>