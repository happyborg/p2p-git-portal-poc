<script>
// import { onMount } from 'svelte';
import wasm from './main.go';
const { uploadFile, getRepositoryList, getHeadCommitsRange, newRepository} = wasm;

import RepoDashboardPanel from './RepoDashboardPanel.svelte'
import IssuesListingPanel from './IssuesListingPanel.svelte'
import CommitsListingPanel from './CommitsListingPanel.svelte'
import DirectoryListingPanel from './DirectoryListingPanel.svelte'

import FileUploadPanel from './test/FileUploadPanel.svelte'
import GoGitClonePanel from './test/GoGitClonePanel.svelte'

let droppedFiles = []
let uploadingFile
let errorMessage

$: manageUploads(droppedFiles)

let newRepoName =''
let activeRepository
let directoryPath = ''

let repositoryPath = ''
$: repositoryPath = (allRepositories && activeRepository !== undefined && allRepositories[activeRepository] !== undefined ?
    allRepositories[activeRepository].path : '')

$: updateDirectoryPath(activeRepository)
$: console.log("directoryPath:", directoryPath)	// Debug

async function updateDirectoryPath(activeRepository) {
	console.log("updateDirectoryPath()", activeRepository)
	let repo = allRepositories[activeRepository]
	if (repo === undefined) {
		directoryPath = ''
		return
	}

	directoryPath = repo.path
}

async function makeNewRepo() {
	newRepoName.trim()
	if (newRepoName.length === 0) {
		errorMessage = "Please enter a name for the new repository"
		return
	}
	else if (getRepositoryForDirectory(newRepoName) !== undefined) {
		errorMessage = "Repository already exists at: " + newRepoName
		return
	}

	await newRepository(newRepoName)
	await updateRepositoryUI(newRepoName)
}

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

async function updateRepositoryUI(newRepoPath) {
	allRepositories = await getRepositoryList();
	setActiveRepository(newRepoPath)
}

function setActiveRepository(repoDirectory) {
	for (let index = 0; index < allRepositories.length; index++) {
		let repo = allRepositories[index]
		if (repo.path === repoDirectory) {
			activeRepository = index
			return
		}
	}
}

function getRepositoryForDirectory(directoryName) {
	for (let index = 0; index < allRepositories.length; index++) {
		if (allRepositories[index].path === directoryName) return repo
	}
	return undefined
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

<p> This is an experimental git portal (like github) that will run entirely in
the browser from static storage, so no server-side code and no third parties involved. Built
using Svelte and Golang/Web Assembly to run on peer-to-peer networks such as <a
href='https://safenetwork.tech'>Safe Network</a>. Read more on github at <a
href='https://github.com/happybeing/p2p-git-portal-poc'>p2p-git-portal-poc</a></p>.

<div class='top-grid'>
	<RepoDashboardPanel bind:activeRepository={activeRepository} bind:allRepositories={allRepositories}></RepoDashboardPanel>
	<IssuesListingPanel bind:repositoryPath={repositoryPath}></IssuesListingPanel>
</div>

<div class='top-grid'>
	<p>
		<button type="button" on:click={() => { makeNewRepo(newRepoName); }}>New Repository:</button>
        <input bind:value={newRepoName} placeholder="directory name"><br/>		
	</p>
	<CommitsListingPanel bind:activeRepository={activeRepository} bind:allRepositories={allRepositories}></CommitsListingPanel>
</div>

{#if errorMessage}
<div>
	<p style="color: #f00">{errorMessage}</p>
	<button type="button" on:click={() => { errorMessage = undefined; }}>Dismiss</button>
</div>
{/if}
<div class='top-grid'>
	<!-- <FileUploadPanel bind:droppedFiles={droppedFiles} bind:errorMessage={errorMessage} >
		<p>Files to upload: {droppedFiles.length}<br/>
			{#if uploadingFile}
				Uploading: {uploadingFile}
			{/if}
			<br/>
		</p>		
	</FileUploadPanel> -->
	<GoGitClonePanel updateRepositoryUI={updateRepositoryUI} bind:errorMessage={errorMessage} ></GoGitClonePanel>
		<DirectoryListingPanel storeName="Storage" bind:directoryPath={directoryPath}></DirectoryListingPanel>
	</div>
<!-- <GoWasmExample bind:errorMessage={errorMessage} ></GoWasmExample> -->
</main>

