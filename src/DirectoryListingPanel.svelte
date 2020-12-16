<script>
// Show a directory in a repository

import wasm from './main.go';
const {getDirectory} = wasm;

export let storeName = "Content"
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
    <h3>{storeName}</h3>
    <h2>
        <!-- <a href="#/" on:click={() => newPath("")}>{storeName}</a> -->
        <a href="#root" on:click={() => newPath("/")}>/</a>
        {#if repositoryRoot.length > 0}
        <a href="#{repositoryRoot}" on:click={() => newPath(repositoryRoot)}>{repositoryRoot}</a>
        {/if}
        {#each crumbs as crumb, i}
            &nbsp;/&nbsp;<a href="#{crumb}" on:click={() => newPath(buildPath(i))}>{crumb}</a>
        {/each}
    </h2>

    {#if directory && directory.length > 0}
        {#each directory as item, index}
            {#if item}
                <span>
                    {#if item.type === "directory"}
                        <a href="#{item.name}" on:click={appendToPath(item.name)}><i>📁</i> {item.name}</a>
                    {:else}
                        <i>📄</i> {item.name}
                    {/if}
                </span><br/>
            {/if}
        {/each}
    {:else}
        No files.
    {/if}
</div>
