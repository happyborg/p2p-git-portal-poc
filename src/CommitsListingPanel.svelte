<script>
// Show commits for the active repository

import wasm from './main.go';
const { getHeadCommitsRange } = wasm;

export let activeRepository;
export let allRepositories;

let repositoryPath = '';
$: repositoryPath = allRepositories && activeRepository !== undefined && allRepositories[activeRepository] !== undefined ?
    allRepositories[activeRepository].path : ''; 

$: updateCommitsListing(repositoryPath);

let commits = [];

async function updateCommitsListing(repoPath) {
    console.log("updateCommitsListing() repoPath: ", repoPath);
    commits = [];
    let result = [];
    let commitsRange;
    if (repoPath) {
        try {
            commitsRange = await getHeadCommitsRange(repoPath, 0, 10);
            console.dir(commitsRange);
            if (commitsRange !== undefined && commitsRange && commitsRange.commits !== undefined) {
                result = commitsRange.commits;
                console.log("result now: ");
                console.dir(result);
            }
        } catch(e) {
            return;
        }
    }

    commits = [...result];
    console.log("commits now: ");
    console.dir(commits);
}

</script>

<style>
/* .highlight {
    font-weight: bold;
} */

</style>
<div>
    <ul>
        <h2>Commit History</h2>
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