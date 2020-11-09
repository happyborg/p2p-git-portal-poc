<script>
import { onMount } from 'svelte';
import wasm from '../main.go';

import NumberInput from './NumberInput.svelte';

const { add, gitClone, raiseError, someValue } = wasm;

let staticValue;

onMount(async () => {
  staticValue = await someValue();
});

let values = [9, 9];
let total;

$: total = add(...values);

function handleChangedValue(event) {
	values[event.detail.index] = event.detail.value;
	values = [...values];
}

export let errorMessage;

async function appRaiseError() {
	try {
		await raiseError();
	} catch (e) {
		errorMessage = e
	}
}
</script>

<hr>
<h2>Go Wasm Example</h2>
<p>Change the values in the boxes below and the total will update.<br/>Click the button to add more input boxes.</p>
{#each values as value, index}
<!-- Version using a Svelte component -->
<NumberInput key={index} {value} on:updateValue={handleChangedValue} />

<!-- Version not using a Svelte component:
<input type="number" value={value} on:change={ (e) => {if (value !== parseInt(e.target.value, 10)) {updateValues(index, parseInt(e.target.value, 10))} }} />
-->
{/each}

<button type="button" on:click={() => {values = [...values, 9]}}>More inputs!</button>

{#await total}
	<p>...calculating</p>
{:then number}
	<p>Total of values is {number}</p>
{:catch error}
	<p style="color: red">{error.message}</p>
{/await}
<div>
	<p>
		Click this button to simulate an error:
		<button type="button" on:click={() => appRaiseError()}>Make error!</button>
	</p>

	{#if errorMessage}
	<div>
		<p style="color: #f00">{errorMessage}</p>
		<button type="button" on:click={() => { errorMessage = undefined; }}>Dismiss</button>
	</div>
	{/if}
</div>
<div>
	<p>Here's a static value: {staticValue}</p>
</div>

