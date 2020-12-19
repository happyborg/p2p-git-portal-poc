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
            }
        } catch(e) {
            return;
        }
    }

    commits = [...result];
    // console.log("commits now: ");
    // console.dir(commits);
}

</script>

<style>
/* .highlight {
    font-weight: bold;
} */

</style>
    <div class="w-full overflow-x-auto">
        <table class="w-full whitespace-no-wrap">
          <tbody
          class="bg-white divide-y dark:divide-gray-700 dark:bg-gray-800">
    
        {#if commits && commits.length > 0}
            {#each commits as commit, index}
                {#if commit}
                <tr class="text-gray-700 dark:text-gray-400">
                    <td class="px-4 py-3">
                      <div class="flex items-center text-sm">
                        <!-- Avatar with inset shadow -->
                        <div
                          class="relative hidden w-8 h-8 mr-3 rounded-full md:block"
                        >
                          <img
                            class="object-cover w-full h-full rounded-full"
                            src="{commit.author_img}"
                            alt=""
                            loading="lazy"
                          />
                          <div
                            class="absolute inset-0 rounded-full shadow-inner"
                            aria-hidden="true"
                          ></div>
                        </div>
                        <div>
                          <p class="font-semibold">{commit.message}</p>
                          <p class="text-xs text-gray-600 dark:text-gray-400">
                            Committed on {commit.date} by {commit.author}
                          </p>
                        </div>
                      </div>
                    </td>
                    <td class="px-4 py-3 text-sm">
                        #{commit.hash.substr(0,7)}
                    </td>
                  </tr>

                {/if}
            {/each}
        {:else}
            No commits.
        {/if}
    </tbody>
</table>
</div>