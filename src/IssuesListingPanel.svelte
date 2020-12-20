<script>
// Show issues for the active repository

import wasm from './main.go';
const { getIssuesForRepo } = wasm;

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
    <ol class="border border-gray-300 rounded bg-white">
        {#if issues && issues.length > 0}
            {#each issues as issue, index}
                {#if issue}
                    <li class="border-t border-gray-300 p-3 first-child:border-t-0">
                    <p class="font-semibold text-sm">{issue.title}</p>
                    <div class="flex items-center text-sm">
                      <div class="text-xs text-gray-600 dark:text-gray-400">
                        #{issue.id.substr(0,7)}
                      </div>
                      </div>
                    </li>
                {/if}
            {/each}
        {:else}
        <li class="border-t border-gray-300 p-3 first-child:border-t-0">
            <p class="font-semibold text-sm">No issues</p>
        </li>
        {/if}
    </ol>
</div>