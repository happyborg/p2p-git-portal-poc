<script>
// import { onMount } from 'svelte';
import wasm from './main.go';
const { uploadFile, getRepositoryList, listHeadCommits, getHeadCommitsRange } = wasm;

import RepoDashboardPanel from './RepoDashboardPanel.svelte'
import CommitsListingPanel from './CommitsListingPanel.svelte'

import FileUploadPanel from './test/FileUploadPanel.svelte'
import GoGitClonePanel from './test/GoGitClonePanel.svelte'

let droppedFiles = [];
let uploadingFile;
let errorMessage;

$: manageUploads(droppedFiles)

let activeRepository = 0;

// Development:
let allRepositories = [];

function readFile(entry, successCallback, errorCallback) {
  entry.file(function(file) {
    let reader = new FileReader();

    reader.onload = function() {
      successCallback(entry, reader.result);
    };

    reader.onerror = function() {
      errorCallback(entry, reader.error);
    }

    reader.readAsArrayBuffer(file);
  }, errorCallback);
}

async function manageUploads(droppedFiles) {
	console.log("manageUploads()");

	let fileInfo;
	while (fileInfo = droppedFiles.pop()) {
		console.log('uploading: ', fileInfo.fullPath);
		// console.dir(fileInfo);
		readFile(fileInfo, 
			async (fileInfo, result) => { 
				console.log('passing to Golang: ', fileInfo.fullPath)
				console.dir(result);
				await uploadFile(...[fileInfo.fullPath, new Uint8Array(result)])
			},
			(fileInfo, result) => { console.log('error loading file: ', fileInfo.fullPath)}
		);
	}
}

async function updateRepositoryUI() {
	allRepositories = await getRepositoryList();
}

async function testRangeCommits() {
	console.log('testRangeCommits()')
	let result = await getHeadCommitsRange("saeedareffard1377666/testproject2.git", 0, 100);
	console.dir(result);
}

async function testReturnTypes() {
	console.log("testReturnTypes()")
	let ret = await testTypes();
	console.dir(ret);
}

</script>

<style>
	main {
		text-align: left;
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;

		font-family: -apple-system,BlinkMacSystemFont,Segoe UI,Helvetica,Arial,sans-serif,Apple Color Emoji,Segoe UI Emoji;
		font-size: 14px;
		line-height: 1.5;
		color: #24292e;
	}

	h1 {
		color: #ff3e00;
		text-transform: uppercase;
		font-size: 2em;
		font-weight: 100;
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}
	.top-grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		/* grid-template-rows: 5fr 20fr 5fr; */
		grid-gap: 10px;
		/* height: 720px; */
	}
</style>

<main>
<h1>p2p Git Portal (POC)</h1>
<div class='top-grid'>
	<RepoDashboardPanel bind:activeRepository={activeRepository} bind:allRepositories={allRepositories}></RepoDashboardPanel>
	<CommitsListingPanel bind:activeRepository={activeRepository} bind:allRepositories={allRepositories}></CommitsListingPanel>
</div>

<div>
	<p>
		<button type="button" on:click={() => { updateRepositoryUI(); }}>Test list repos</button><br/>
		<button type="button" on:click={() => { listHeadCommits("http://localhost:8010/proxy/happybeing/p2p-git-portal-poc.git"); }}>Test list HEAD commits</button><br/>
		<button type="button" on:click={() => { testRangeCommits() }}>Test get range HEAD commits</button><br/>
	</p>
</div>

{#if errorMessage}
<div>
	<p style="color: #f00">{errorMessage}</p>
	<button type="button" on:click={() => { errorMessage = undefined; }}>Dismiss</button>
</div>
{/if}
<div class='top-grid'>
	<FileUploadPanel bind:droppedFiles={droppedFiles} bind:errorMessage={errorMessage} >
		<p>Files to upload: {droppedFiles.length}<br/>
			{#if uploadingFile}
				Uploading: {uploadingFile}
			{/if}
			<br/>
		</p>		
	</FileUploadPanel>
	<GoGitClonePanel updateRepositoryUI={updateRepositoryUI} bind:errorMessage={errorMessage} ></GoGitClonePanel>
</div>

<!-- <GoWasmExample bind:errorMessage={errorMessage} ></GoWasmExample> -->
</main>

