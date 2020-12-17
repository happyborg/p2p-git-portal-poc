<script>
// Show issues for the active repository

import wasm from './main.go';
const { getIssuesForRepo } = wasm;

export let disabled = true
export let repositoryRoot;

$: updateIssuesListing(repositoryRoot);

let issues = [];

async function updateIssuesListing(repoPath) {
    console.log("updateIssuesListing() repoPath: ", repoPath);
    let result = [];
    if (repoPath) {
        result = await getIssuesForRepo(repoPath);
    }

    issues = [...result];
    // console.log("issues now: ");
    // console.dir(issues);
}

</script>

<style>
/* .highlight {
    font-weight: bold;
} */

</style>
<div>
    <h3>Issues</h3>
    {#if issues && issues.length > 0}
        {#each issues as issue, index}
            {#if issue}
                <span>#{issue.id.substr(0,7)} {issue.title}</span><br/>
            {/if}
        {/each}
    {:else}
        No issues.
    {/if}
</div>