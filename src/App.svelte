<script>
// import { onMount } from 'svelte';
import wasm from './main.go';
const { wasmReady, uploadFile, getRepositoryList, getHeadCommitsRange, newRepository, openRepository} = wasm;

import RepoDashboardPanel from './RepoDashboardPanel.svelte'
import IssuesListingPanel from './IssuesListingPanel.svelte'
import CommitsListingPanel from './CommitsListingPanel.svelte'
import DirectoryListingPanel from './DirectoryListingPanel.svelte'

import UploadRepositoryPanel from './test/UploadRepositoryPanel.svelte'
import GoGitClonePanel from './test/GoGitClonePanel.svelte'

import Tab from './components/tabs/Tab.svelte'
import Tabs from './components/tabs/Tabs.svelte'
import TabList from './components/tabs/TabList.svelte'
import TabPanel from './components/tabs/TabPanel.svelte'
import Alert from './components/alert/Alert.svelte'

let uploadRoot = ''
let filesToUpload = []
let uploadStatus = ''
let currentUpload = ''
let errorMessage

$: manageUploads(filesToUpload)

let newRepoName =''
let activeRepository
let directoryPath = ''

let repositoryRoot = ''
$: repositoryRoot = (allRepositories && activeRepository !== undefined && allRepositories[activeRepository] !== undefined ?
	allRepositories[activeRepository].path.trim() : '')

$: updateDirectoryPath(activeRepository)
$: console.log("directoryPath:", directoryPath)	// Debug

let appLoading = true
async function waitUntilReady() {
	errorMessage = "Loading application (wasm) - please WAIT until this message disappears"
	appLoading = !await wasmReady()
	errorMessage = undefined
}
waitUntilReady();

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

async function readFile(entry, successCallback, errorCallback) {
    let reader = new FileReader();

    reader.onload = async function() {
      successCallback(entry, reader.result);
    };

    reader.onerror = function() {
      errorCallback(entry, reader.error);
    }

    reader.readAsArrayBuffer(entry);
}

async function manageUploads(filesToUpload) {
	console.log("manageUploads()");
	if (filesToUpload === undefined || filesToUpload.length === 0) return

	// TODO: check if repository exists
	await uploadFiles(uploadRoot, filesToUpload)
}

async function uploadFiles(uploadRoot, filesToUpload) {
	console.log("uploadFiles() to", uploadRoot);
	const totalFiles = filesToUpload.length
	let filesUploaded = 0
	uploadStatus = "Uploading " + totalFiles + " files to " + uploadRoot + "/"

	let fileInfo;
	while (fileInfo = filesToUpload.pop()) {
		const fullPath = fileInfo.webkitRelativePath
		// console.log('uploading: ', fileInfo.webkitRelativePath);
		currentUpload = fileInfo.webkitRelativePath.substring(uploadRoot+1)
		// console.dir(fileInfo)
		await readFile(fileInfo, 
			async (fileInfo, result) => { 
				// console.log('passing to Golang: ', fullPath)
				// console.dir(result);
				await uploadFile(...[fullPath, new Uint8Array(result)])
				filesUploaded += 1
				uploadStatus = filesUploaded + " files uploaded to " + uploadRoot + "/"
				if (filesUploaded === totalFiles) {
					// TODO: open repository and select it
					await openRepository(uploadRoot, (error) => {
						if (error && error !== "") {
							let message = "Failed to open repository of uploaded directory: " + error;
							errorMessage = message;
							console.log(message);
						} else {
							console.log("Upload complete and repository opened:", uploadRoot);
							updateRepositoryUI(uploadRoot)
						}
					});
				}
			},
			(fileInfo, result) => { console.log('error loading file: ', fullPath)}
		);
	}

	currentUpload = ''
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
		if (allRepositories[index].path === directoryName) return allRepositories[index].repo
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

<style global lang="postcss">
	/* only apply purgecss on utilities, per Tailwind docs */
	/* purgecss start ignore */
	@tailwind base;
	@tailwind components;
	/* purgecss end ignore */
	@tailwind utilities;
</style>
<nav class="flex items-center z-20 h-14 relative text-gray-700 border-b-2 ">
	<a class="text-l font-medium leading-none flex place-items-center pl-6 py-5" href="/">
		<svg xmlns="http://www.w3.org/2000/svg" 
			class="w-9 h-9 mr-2"
			fill="#f0572a"
			viewBox="0 0 24.912 28.218"><path d="M23.893 18.65a1.604 1.604 0 00-.679-.17 1.731 1.731 0 00-.972.22l-.154.083-1.25-.723-1.252-.723-.012-.237a1.598 1.598 0 00-.154-.654 1.859 1.859 0 00-.793-.847l-.069-.035v-2.898l.125-.073.125-.075 1.25.723 1.252.723.012.224a1.79 1.79 0 001.042 1.557c.251.118.442.16.752.16.332 0 .511-.044.81-.193.455-.226.783-.627.928-1.132a2.164 2.164 0 000-.934 1.795 1.795 0 00-.928-1.132 1.527 1.527 0 00-.81-.191c-.334 0-.515.042-.816.193l-.214.108-.502-.291-1.256-.725-.754-.434.006-.143.006-.143 1.252-.723 1.25-.723.154.083c.32.172.611.239.972.22a1.752 1.752 0 001.223-.571c.334-.363.492-.801.469-1.3a1.481 1.481 0 00-.18-.708c-.451-.926-1.564-1.275-2.486-.776l-.151.081-1.25-.723-1.252-.723-.012-.237a1.598 1.598 0 00-.154-.654 1.78 1.78 0 00-1.111-.974 2.018 2.018 0 00-.885-.042A1.79 1.79 0 0015.992 4.8c-.014.253.021.465.121.727.164.438.54.828.968 1.005.266.11.403.135.72.135.255-.002.313-.008.478-.058.104-.031.266-.098.359-.149l.175-.093.336.193 1.254.723.918.529v.29l-.9.518c-.492.286-1.057.61-1.25.723l-.355.203-.205-.104a1.752 1.752 0 00-1.42-.1 2.161 2.161 0 00-.301.135c-.1.058-.133.069-.166.05a92.42 92.42 0 01-1.254-.723l-1.21-.702-.01-.207a1.804 1.804 0 00-.33-.974 2.135 2.135 0 00-.55-.509l-.147-.087V4.878l.002-1.447.154-.093a1.798 1.798 0 00.756-2.18 1.8 1.8 0 00-.714-.87 1.858 1.858 0 00-1.204-.272 2.8 2.8 0 00-.561.17 1.808 1.808 0 00-.953 2.01c.106.475.411.895.828 1.142l.154.093.002 1.447v1.447l-.147.087a1.93 1.93 0 00-.714.789 1.57 1.57 0 00-.162.756c0 .332.044.513.193.81.229.455.627.785 1.121.924.436.125.988.054 1.379-.179l.12-.071 1.254.72 1.254.723-.006.152-.004.152-1.225.708c-.673.388-1.24.712-1.256.716-.019.004-.087-.025-.156-.066a1.868 1.868 0 00-1.071-.22 2.05 2.05 0 00-.691.22.895.895 0 01-.141.073c-.01 0-.579-.324-1.264-.718l-1.244-.718v-.135c0-.305-.095-.656-.255-.93a2.308 2.308 0 00-.624-.622l-.141-.081V6.5l.141-.081c.189-.108.513-.432.621-.621.106-.181.202-.451.233-.648a1.796 1.796 0 00-.494-1.524 1.788 1.788 0 00-.924-.507 2.018 2.018 0 00-.885.042 1.785 1.785 0 00-1.113.978c-.12.26-.16.469-.149.808.008.26.017.314.079.492.143.428.465.816.839 1.013l.095.052V9.41l-.127.069-.124.071-.357-.206c-.198-.114-.762-.438-1.254-.724l-.9-.517v-.158a1.8 1.8 0 00-1.795-1.784c-.688 0-1.308.388-1.609 1.007-.15.308-.177.422-.177.789 0 .365.023.469.176.789.156.324.436.625.719.773l.112.058v2.912l-.139.081a2.22 2.22 0 00-.59.567 1.724 1.724 0 00-.282 1.059c.025.517.239.955.633 1.296.166.143.488.309.72.37.235.062.639.062.872 0a1.802 1.802 0 001.233-2.411 1.948 1.948 0 00-.751-.88l-.139-.081V9.577l.125-.067.123-.066 1.252.723 1.252.723.014.247c.019.363.118.654.316.932.125.172.39.417.558.517l.139.081v2.896l-.148.087c-.289.17-.571.486-.712.793a1.594 1.594 0 00-.156.658l-.012.237-1.252.723-1.25.723-.172-.091a2.488 2.488 0 00-.347-.147c-.162-.054-.203-.058-.498-.058-.363-.002-.467.019-.764.156a1.796 1.796 0 00-1.013 1.283c-.042.21-.031.606.023.801.052.195.224.525.351.679.233.278.567.498.916.6.191.054.542.075.747.042a1.805 1.805 0 001.526-1.765v-.158l.899-.517c.492-.286 1.057-.61 1.252-.723l.353-.203.207.102c.301.152.476.193.818.195.395 0 .689-.083.995-.287a1.8 1.8 0 00-.091-3.057l-.114-.06V12.65l.12-.064.12-.064 1.254.725 1.256.725.004.201c.004.206.046.405.131.631.112.299.434.671.733.847l.158.093.002 1.445v1.447l-.149.087a2.264 2.264 0 00-.561.525 1.82 1.82 0 00-.318.959l-.006.208-1.252.725c-1.129.652-1.256.72-1.295.696a1.807 1.807 0 00-2.689.953c-.062.179-.071.233-.079.492-.006.245-.002.32.035.473.079.332.243.627.484.866.22.218.434.347.743.446.237.077.639.096.885.042a1.79 1.79 0 001.256-.978c.106-.22.164-.448.179-.714l.012-.218 1.25-.723 1.252-.721.127.075.127.075v1.447l-.002 1.445-.155.094c-.6.355-.955 1.073-.862 1.742.077.565.353 1.003.82 1.308.293.191.592.278.966.278s.673-.087.965-.278c.37-.241.619-.565.743-.966.06-.191.104-.473.093-.604l-.006-.087 1.256-.727 1.256-.725.062.037c.156.096.343.174.527.222.178.046.23.05.49.042.312-.01.513-.058.748-.177a1.82 1.82 0 00.947-1.202c.037-.154.042-.228.035-.473-.008-.26-.017-.314-.079-.492a1.815 1.815 0 00-1.219-1.173c-.16-.048-.22-.054-.473-.056-.318 0-.455.025-.72.135a1.893 1.893 0 00-.8.671 1.929 1.929 0 00-.297.98v.15l-1.25.72-1.25.723-.131-.071-.131-.071V21.904l.149-.087c.176-.106.448-.359.561-.525.195-.289.303-.615.318-.959l.008-.208 1.21-.702c.667-.386 1.229-.712 1.254-.723.033-.019.066-.008.166.05.066.039.204.1.301.135.475.168.955.135 1.418-.1l.208-.104.353.203c.195.112.76.436 1.252.723l.899.517v.158c0 .642.378 1.271.937 1.566.315.166.616.23.969.212a1.787 1.787 0 001.686-1.79 1.778 1.778 0 00-1.016-1.62zm-9.647-4.478l.004-.201 1.256-.727c1.198-.691 1.258-.725 1.304-.696.027.017.087.052.135.077l.087.048-.006 1.441-.006 1.443-.16.093c-.206.12-.486.397-.604.596a1.942 1.942 0 00-.272.959l-.002.137-1.252.723-1.254.72-.127-.075-.126-.075v-1.447l.002-1.445.158-.093c.51-.301.853-.884.863-1.478zm4.346 3.814a1.02 1.02 0 01-.791.322c-.313 0-.542-.089-.768-.295l-.131-.12-1.68.968c-.922.534-1.686.976-1.694.984-.01.008-.002.071.014.141.044.162.042.442 0 .588-.114.382-.434.704-.791.795l-.1.025-.006 1.952-.004 1.95.156.05c.278.085.536.307.665.573.081.164.114.311.114.511 0 .316-.108.569-.341.797a1.02 1.02 0 01-.791.322c-.102 0-.232-.012-.291-.029a1.174 1.174 0 01-.82-.858 1.349 1.349 0 01.046-.631c.129-.349.44-.627.791-.708l.097-.023v-3.896l-.087-.023c-.417-.106-.764-.459-.847-.866a1.349 1.349 0 01.046-.631c.129-.349.44-.627.791-.708l.097-.023v-3.896l-.087-.023c-.417-.106-.764-.459-.847-.866a1.44 1.44 0 01.042-.625 1.18 1.18 0 01.708-.691c.212-.069.527-.069.725 0 .357.122.64.428.737.793.042.152.039.44-.002.581-.114.382-.434.704-.791.795l-.1.025-.006 1.952-.004 1.95.154.05c.15.046.33.152.459.272l.062.056 1.705-.984 1.705-.984-.029-.108a1.434 1.434 0 01.035-.654 1.17 1.17 0 01.716-.691l.156-.05v-3.87l-.087-.023a1.162 1.162 0 01-.752-.602.998.998 0 01-.11-.571 1.12 1.12 0 01.866-1.04c.16-.039.453-.025.623.031.355.118.639.426.737.793.042.152.04.44-.002.581-.114.382-.434.704-.791.795l-.1.025-.006 1.945c-.004 1.844-.002 1.945.031 1.945.085.002.326.106.45.195.303.216.469.542.469.928-.001.317-.109.57-.341.799z"/><path d="M17.668 10.503a.562.562 0 00-.307.928c.131.141.222.181.419.181.197 0 .289-.039.419-.183a.57.57 0 00-.029-.787.606.606 0 00-.502-.139zM12.332 13.576a.562.562 0 00-.307.928c.131.141.222.181.419.181s.289-.039.419-.183a.57.57 0 00-.029-.787.602.602 0 00-.502-.139zM17.689 16.628a.562.562 0 00-.307.928c.131.141.222.181.419.181.197 0 .289-.039.419-.183a.57.57 0 00-.029-.787.606.606 0 00-.502-.139zM12.332 19.722a.562.562 0 00-.307.928c.131.141.222.181.419.181s.289-.039.419-.183a.57.57 0 00-.029-.787.602.602 0 00-.502-.139zM12.332 25.868a.562.562 0 00-.307.928c.131.141.222.181.419.181s.289-.039.419-.183a.57.57 0 00-.029-.787.599.599 0 00-.502-.139z"/>
		</svg> Git Portal
	</a>
	
	<a class="leading-none flex ml-auto place-items-center text-sm inline-block bg-red text-gray-700 h-full px-6" href="https://github.com/happyborg/p2p-git-portal-poc">
		<svg xmlns="http://www.w3.org/2000/svg" 
			class="w-5 h-5 mr-2"
			fill="currentColor"
			viewBox="0 0 32.58 31.77">
			<path d="M16.37.06C7.38.06.08 7.35.08 16.35c0 7.2 4.67 13.3 11.14 15.46.81.15 1.11-.35 1.11-.79 0-.39-.01-1.41-.02-2.77-4.53.98-5.49-2.18-5.49-2.18-.74-1.88-1.81-2.38-1.81-2.38-1.48-1.01.11-.99.11-.99 1.63.12 2.5 1.68 2.5 1.68 1.45 2.49 3.81 1.77 4.74 1.35.15-1.05.57-1.77 1.03-2.18-3.62-.41-7.42-1.81-7.42-8.05 0-1.78.63-3.23 1.68-4.37-.17-.41-.73-2.07.16-4.31 0 0 1.37-.44 4.48 1.67 1.3-.36 2.69-.54 4.08-.55 1.38.01 2.78.19 4.08.55 3.11-2.11 4.48-1.67 4.48-1.67.89 2.24.33 3.9.16 4.31 1.04 1.14 1.67 2.59 1.67 4.37 0 6.26-3.81 7.63-7.44 8.04.58.5 1.11 1.5 1.11 3.02 0 2.18-.02 3.93-.02 4.47 0 .44.29.94 1.12.78 6.47-2.16 11.13-8.26 11.13-15.45 0-9-7.3-16.3-16.29-16.3z" />
		</svg>  View on Github
	</a>
</nav>
<div class="flex flex-wrap bg-gray-50 dark:bg-gray-900 h-full text-gray-700 dark:text-gray-400">
	<aside class="md:left-0 md:block md:fixed md:top-0 md:bottom-0 md:overflow-y-auto md:flex-row md:flex-no-wrap md:overflow-hidden border-solid border-r border-light-gray bg-white flex flex-wrap items-center justify-between relative md:w-96 z-10 pb-4 px-6 pt-14">
		<RepoDashboardPanel bind:activeRepository={activeRepository} bind:allRepositories={allRepositories}></RepoDashboardPanel>
		
		<GoGitClonePanel updateRepositoryUI={updateRepositoryUI} bind:errorMessage={errorMessage} ></GoGitClonePanel>
		<div class="my-3">
			<label class="block text-gray-800 text-m font-semibold mb-2 mt-6 pt-4 border-t" for="">Create repository</label>
			<input bind:value={newRepoName} class="appearance-none w-full p-3 text-xs font-semibold leading-none bg-gray-100 rounded outline-none" type="text" name="field-name" placeholder="Repository name">
		  </div>
		<p>
			<button  disabled={appLoading} type="button" on:click={() => { makeNewRepo(newRepoName); }} class="bg-white text-sm text-gray-800 font-bold border border-gray-200 hover:border-red-600 hover:bg-red-500 hover:text-white py-2 px-5 inline-flex items-center">
				<span class="mr-2">Create</span>
				<svg fill="currentColor" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M2 2.5A2.5 2.5 0 014.5 0h8.75a.75.75 0 01.75.75v12.5a.75.75 0 01-.75.75h-2.5a.75.75 0 110-1.5h1.75v-2h-8a1 1 0 00-.714 1.7.75.75 0 01-1.072 1.05A2.495 2.495 0 012 11.5v-9zm10.5-1V9h-8c-.356 0-.694.074-1 .208V2.5a1 1 0 011-1h8zM5 12.25v3.25a.25.25 0 00.4.2l1.45-1.087a.25.25 0 01.3 0L8.6 15.7a.25.25 0 00.4-.2v-3.25a.25.25 0 00-.25-.25h-3.5a.25.25 0 00-.25.25z"></path></svg>
			</button>
		</p>
		<UploadRepositoryPanel bind:disabled={appLoading} bind:uploadRoot={uploadRoot} bind:filesToUpload={filesToUpload}></UploadRepositoryPanel>
		{#if uploadStatus !== ''}
			{uploadStatus}<br/>
			{currentUpload}
		{/if}
	</aside>
		<main class="h-screen -mt-14 pt-14 overflow-y-auto w-full px-4 pb-4 md:ml-96 bg-gray-50">
			<Alert type="info">
			This is an experimental git portal (like github) that will run entirely in
				the browser from static storage, so no server-side code and no third parties involved. Built
				using Svelte and Golang/Web Assembly to run on peer-to-peer networks such as <a class="underline"
				href='https://safenetwork.tech'>Safe Network</a>.
			</Alert>
			{#if errorMessage}
				<Alert type="error" dismissable="true" bind:errorMessage={errorMessage}>
					<p>{errorMessage}</p>
				</Alert>
			{/if}
			{#if repositoryRoot}
				<Tabs>
					<TabList>
						<Tab>Code</Tab>
						<Tab>Commits</Tab>
						<Tab>Issues</Tab>
					</TabList>

					<TabPanel>
						<!-- <DirectoryListingPanel storeName="Storage" bind:directoryPath={directoryPath}></DirectoryListingPanel> -->
						<DirectoryListingPanel bind:disabled={appLoading} storeName="Worktree" bind:repositoryRoot={repositoryRoot}></DirectoryListingPanel>
					</TabPanel>
					<TabPanel>
						<CommitsListingPanel bind:disabled={appLoading} bind:repositoryRoot={repositoryRoot}></CommitsListingPanel>
					</TabPanel>
					<TabPanel>
						<IssuesListingPanel bind:disabled={appLoading} bind:repositoryRoot={repositoryRoot}></IssuesListingPanel>
					</TabPanel>
				</Tabs>
			{/if}


</main>

</div>