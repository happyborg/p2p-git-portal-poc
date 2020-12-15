<script>
// Select active repository from those stored in the Go filesystem

export let activeRepository;
export let allRepositories;

</script>

<style>
.highlight {
    font-weight: bold;
}

</style>
<div>
    <h2>Repositories</h2>
    {#if allRepositories !== undefined && allRepositories.length > 0}
        {#each allRepositories as repo, index}
            <span class={index == activeRepository ? 'highlight' : ''}>
                <a href= '#none' on:click={() => activeRepository = index}>{repo.path}</a></span><br/>
        {/each}
    {:else}
        No repositories found.<br/>
        Add repositories using New, Upload or Clone.
    {/if}
</div>