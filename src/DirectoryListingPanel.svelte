<script>
import { utimes } from 'fs';

// Show a directory in a repository

import wasm from './main.go';
import FolderIcon from './components/icons/Folder.svelte'
import FileIcon from './components/icons/File.svelte'
const {getDirectory} = wasm;

//export let storeName = "Content"
export let repositoryRoot = ''

let directoryPath = ''
let directory = []
let repositoryCrumbs = []
let crumbs = []

$: updateRepository(repositoryRoot)
$: updateCrumbs(directoryPath)
$: updateDirectory(directoryPath)

async function updateRepository(repositoryRoot) {
    // console.log("updateRepository(): ", repositoryRoot)
    repositoryCrumbs = repositoryRoot.split('/')
    console.dir(repositoryCrumbs)
    directoryPath = repositoryRoot
}

async function updateDirectory(directoryPath) {
    // console.log("updateDirectory(): ", directoryPath)    
	let listing = await getDirectory(directoryPath)
    directory = [...listing]
    directory.sort((a, b) => (a.name.toLowerCase() > b.name.toLowerCase()) ? 1 : -1)
	// console.dir(directory)

	for (let i = 0; i < directory.length; i++) {
		let entry = directory[i]
		// console.log("type: ", entry["type"]);
		// console.log("name: ", entry["name"]);
		// console.log("size: ", Number(entry["size"]));
		// console.log("modified: ", Number(entry["modified"]));
		// console.log("mode: ", entry["mode"]);
		// console.log("sys: ", Number(entry["sys"]));
	}
}

async function updateCrumbs(directoryPath) {
    console.log("updateCrumbs()", directoryPath)
    let newCrumbs = [];
    if (directoryPath !== undefined && directoryPath !== "") {
        let crumbPart = directoryPath.substring(repositoryRoot.length)
        newCrumbs = crumbPart.split('/')
        if (newCrumbs.length > 0 && newCrumbs[0].length === 0) newCrumbs.shift()
    }
    crumbs = [...newCrumbs]
    // console.log("crumbs:")
    // console.dir(crumbs)

    await updateDirectory(directoryPath)
}

function buildPath(n) {
    // console.log("buildPath()", n)
    // console.log("crumbs:")
    // console.dir(crumbs)
    let path = [...repositoryCrumbs, ...crumbs.slice(0, n+1)].join('/')
    console.log("buildPath returns: ", path)
    return path
}

async function newPath(path) {
    // console.log("newPath(): ", path)
    directoryPath = await path
    await updateDirectory(directoryPath)
}

async function appendToPath(itemName) {
    if ( directoryPath.length > 0 ) directoryPath += "/"
    directoryPath += itemName
}

</script>

<style>
/* .highlight {
    font-weight: bold;
} */

</style>
<div>
  
    <div class="mb-2 text-sm">
        <!-- <a href="#/" on:click={() => newPath("")}>{storeName}</a> -->
        <a href="#root" on:click={() => newPath("/")}>&nbsp;/&nbsp;</a>
        {#if repositoryRoot.length > 0}
        <a href="#{repositoryRoot}" on:click={() => newPath(repositoryRoot)}> {repositoryRoot} </a>
        {/if}
        {#each crumbs as crumb, i}
            &nbsp;/&nbsp;<a href="#{crumb}" on:click={() => newPath(buildPath(i))}> {crumb} </a>
        {/each}
    </div>
    
    {#if crumbs.length == 1}
        <div class="bg-gray-100 px-4 pb-2 pt-1 border border-color-gray-300 rounded-t-sm border-b-0">
            <a href="#{repositoryRoot}" on:click={() => newPath(repositoryRoot)}> ...</a>
        </div>
    {/if}

    {#if crumbs.length > 1}
        <div class="bg-gray-100 px-4 pb-2 pt-1 border border-color-gray-300 rounded-t-sm border-b-0">
            <a href="#{crumbs[crumbs.length-2]}" on:click={() => newPath(buildPath(crumbs.length - 2))}> ...</a>
        </div>
    {/if}
    {#if directory && directory.length > 0}
        <div class="border border-color-gray-300 rounded-b-sm text-sm border-t-0 bg-white">
            {#each directory as item, index}
                {#if item && item.type === "directory"}
                <div class="border-t border-color-gray-300 px-4 py-2">
                    <span>
                        <a href="#{item.name}" on:click={appendToPath(item.name)}><FolderIcon cssClass="text-blue-400" /> {item.name}</a>
                    </span>
                </div>
                    {/if}
                
            {/each}
            {#each directory as item, index}
                {#if item && item.type === "file"}
                <div class="border-t border-color-gray-300 px-4 py-2">
                    <span>
                        <FileIcon /> {item.name}
                    </span>
                </div>
                {/if}
            {/each}
        </div>
    {:else}
        No files.
    {/if}
</div>