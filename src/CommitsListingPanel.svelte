<script>
// Show commits for the active repository

export let activeRepository;
export let allRepositories;

let repositoryPath = '';
$: repositoryPath = allRepositories && activeRepository !== undefined && allRepositories[activeRepository] !== undefined ?
    allRepositories[activeRepository].path : 'no repository selected'; 

// Development:
let commits
 = [
    {hash: "d1f23d5e", message: "Initial commit"},
    {hash: "1eb2b4a6", message: "Add some stuff"},
    {hash: "e12c325f", message: "Add lots more stuff"},
];

</script>

<style>
.highlight {
    font-weight: bold;
}

</style>
<div>
    <ul>
        <h2>Commit History for: {repositoryPath}</h2>
        {#if commits && commits.length > 0}
            {#each commits as commit, index}
               <span>#{commit.hash} {commit.message}</span><br/>
            {/each}
        {:else}
            No commits.
        {/if}
    </ul>
</div>