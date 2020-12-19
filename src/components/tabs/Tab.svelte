<script>
	import { getContext } from 'svelte';
	import { TABS } from './Tabs.svelte';

	const tab = {};
	const { registerTab, selectTab, selectedTab } = getContext(TABS);

	registerTab(tab);
</script>

<style>
	button {
		background: none;
		border: none;
 	    border-radius: 0;
		margin: 0;
		color: #ccc;
	}
	
	.selected {
		border-bottom: 2px solid #f0572a;
		color: #333;
	}
</style>

<button class:selected="{$selectedTab === tab}" class="text-gray-600 py-4 px-6 focus:outline-none border-b-2 font-medium border-solid border-light-gray" on:click="{() => selectTab(tab)}">
	<slot></slot>
</button>