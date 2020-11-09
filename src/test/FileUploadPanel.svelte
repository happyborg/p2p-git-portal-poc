<script>
// import { onMount } from 'svelte';
// import wasm from '../main.go';
// const { gitClone, raiseError } = wasm;

import Dropzone from "svelte-file-dropzone";

// export let errorMessage;

let files = {
    accepted: [],
    rejected: []
};

function handleFilesSelect(e) {
    const { acceptedFiles, fileRejections } = e.detail;
    files.accepted = [...files.accepted, ...acceptedFiles];
    files.rejected = [...files.rejected, ...fileRejections];
}
</script>

<div>
    <hr><h2>Upload Files</h2>
    <Dropzone on:drop={handleFilesSelect} />
    <ol>
        {#each files.accepted as item}
        <li>{item.name}</li>
        {/each}
    </ol>
</div>

