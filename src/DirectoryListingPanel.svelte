<script>
// Show a directory in a repository

import wasm from './main.go';
const {getDirectory} = wasm;

export let storeName = "storage"
export let directoryPath;

let directory = [];
let crumbs = [];

$: updateCrumbs(directoryPath)
$: updateDirectory(directoryPath)

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
    // console.log("updateCrumbs()", directoryPath)
    let newCrumbs = [];
    if (directoryPath !== undefined && directoryPath !== "") {
        newCrumbs = directoryPath.split('/')
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
    let path = crumbs.slice(0, n+1).join('/')
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
    <h2>
        <a href="#/" on:click={() => newPath("")}>{storeName}</a>
        <b> / </b>
        {#each crumbs as crumb, i}
            <a href="#{crumb}" on:click={() => newPath(buildPath(i))}>{crumb}</a>
            <b> / </b>
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