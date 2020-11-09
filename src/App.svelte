<script>
// import { onMount } from 'svelte';
// import wasm from './main.go';
// const { add, gitClone, raiseError, someValue } = wasm;

import FileUploadPanel from './test/FileUploadPanel.svelte'
import GoGitClonePanel from './test/GoGitClonePanel.svelte'
import GoWasmExample from './test/GoWasmExample.svelte';

let errorMessage;

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
{#if errorMessage}
<div>
	<p style="color: #f00">{errorMessage}</p>
	<button type="button" on:click={() => { errorMessage = undefined; }}>Dismiss</button>
</div>
{/if}

<div class='top-grid'>
	<FileUploadPanel bind:errorMessage={errorMessage} ></FileUploadPanel>
	<GoGitClonePanel bind:errorMessage={errorMessage} ></GoGitClonePanel>
</div>

<GoWasmExample bind:errorMessage={errorMessage} ></GoWasmExample>
</main>

