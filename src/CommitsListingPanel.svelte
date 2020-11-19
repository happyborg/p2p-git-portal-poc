<script>
// Show commits for the active repository

import wasm from './main.go';
const { getHeadCommitsRange } = wasm;

export let activeRepository;
export let allRepositories;

let repositoryPath = '';
$: repositoryPath = allRepositories && activeRepository !== undefined && allRepositories[activeRepository] !== undefined ?
    allRepositories[activeRepository].path : 'no repository selected'; 

$: updateCommitsListing(repositoryPath);

// Development:
let commits = [];
// commits = [
//     {hash: "d1f23d5e", message: "Initial commit"},
//     {hash: "1eb2b4a6", message: "Add some stuff"},
//     {hash: "e12c325f", message: "Add lots more stuff"},
// ];

async function updateCommitsListing(repoPath) {
    console.log("updateCommitsListing() repoPath: ", repoPath);
    let result = [];
    if (repoPath) {
        let commitsRange = await getHeadCommitsRange(repoPath, 0, 10);
        console.dir(commitsRange);
        if (commitsRange) {
            result = commitsRange.commits;
            console.log("result now: ");
            console.dir(result);
        }
    }

    commits = [...result];
    console.log("commits now: ");
    console.dir(commits);
}

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
                {#if commit}
                    <span>#{commit.hash.substr(0,7)} {commit.message}</span><br/>
                {/if}
            {/each}
        {:else}
            No commits.
        {/if}
    </ul>
</div>