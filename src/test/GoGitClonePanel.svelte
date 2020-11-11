<script>
// import { onMount } from 'svelte';
import wasm from '../main.go';

const { gitClone, testGitCloneDeadlock, listFiles } = wasm;

export let errorMessage;

async function doGitClone() {
	try {
		await testGitCloneDeadlock();
	} catch (e) {
		errorMessage = e
	}
}

</script>
<div>
	<hr><h2>Go Git Clone Panel</h2>
	<p>
		<button type="button" on:click={() => { doGitClone(); }}>Test git clone</button>
	</p>
</div>