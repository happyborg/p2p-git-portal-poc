<script>
// Show issues for the active repository

import wasm from './main.go';
const { getIssuesForRepo } = wasm;

export let activeRepository;
export let allRepositories;

let repositoryPath = '';
$: repositoryPath = allRepositories && activeRepository !== undefined && allRepositories[activeRepository] !== undefined ?
    allRepositories[activeRepository].path : ''; 

$: updateIssuesListing(repositoryPath);

// Development:
let issues = [];
// issues = [
//     {id: "d1f23d5e", title: "Initial commit"},
//     {id: "1eb2b4a6", title: "Add some stuff"},
//     {id: "e12c325f", titletitle: "Add lots more stuff"},
// ];

async function updateIssuesListing(repoPath) {
    console.log("updateIssuesListing() repoPath: ", repoPath);
    let result = [];
    if (repoPath) {
        result = await getIssuesForRepo(repoPath);
        console.log("result now: ");
        console.dir(result);
    }

    issues = [...result];
    console.log("issues now: ");
    console.dir(issues);
}

</script>

<style>
/* .highlight {
    font-weight: bold;
} */

</style>
<div>
    <ul>
        <h2>Issues</h2>
        {#if issues && issues.length > 0}
            {#each issues as issue, index}
                {#if issue}
                    <span>#{issue.id.substr(0,7)} {issue.title}</span><br/>
                {/if}
            {/each}
        {:else}
            No issues.
        {/if}
    </ul>
</div>