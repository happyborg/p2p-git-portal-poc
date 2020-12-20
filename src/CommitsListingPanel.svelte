<script>
// Show commits for the active repository

import wasm from './main.go';
import CommitIcon from './components/icons/Commit.svelte'

const { getHeadCommitsRange } = wasm;

export let repositoryRoot;

$: updateCommitsListing(repositoryRoot);

let commits = [];
let dateGroupArray = [];

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
    commits.sort((a, b) => b.timestamp - a.timestamp)
    const dateGroups = commits.reduce((groups, commit) => {
      const date = commit.date;
      if (!groups[date]) {
        groups[date] = [];
      }
      groups[date].push(commit);
      return groups;
    }, {});

    dateGroupArray = Object.keys(dateGroups).map((date) => {
      return {
        date,
        commits: dateGroups[date]
      };
    });

    console.log("commits now: ");
    console.dir(commits);
}


</script>

<style>
/* .highlight {
    font-weight: bold;
} */

</style>
    <div class="w-full overflow-x-auto relative pl-12">
      <div class="border-l-2 border-gray-300 absolute h-full top-0 left-5"></div>
       {#if dateGroupArray && dateGroupArray.length > 0}
            {#each dateGroupArray as date, index}
                <div class="text-sm my-3">
                  <div style="margin-left: -1px" class="absolute left-3 rounded-full h-5 w-5 bg-gray-50"><CommitIcon cssClass="text-gray-600" /></div>
                   Commits on {date.date}
                </div>
              <ol class="border border-gray-300 rounded bg-white">
                {#each date.commits as commit, index}
                  {#if commit}
                    <li class="border-t border-gray-300 p-3 first-child:border-t-0">
                      <p class="font-semibold text-sm">{commit.message}</p>
                      <div class="flex items-center text-sm">
                        <div class="relative hidden w-6 h-6 mr-3 rounded-full md:block">
                          <img class="object-cover w-full h-full rounded-full" src="{commit.author_img}" alt="" />
                          <div class="absolute inset-0 rounded-full shadow-inner" aria-hidden="true"></div>
                        </div>
                        <div class="text-xs text-gray-600 dark:text-gray-400">
                          Committed on {commit.date} by {commit.author}
                        </div>
                        <div class="ml-auto">
                          <span>#{commit.hash.substr(0,7)}</span>
                        </div>
                      </div>
                    </li>
                  {/if}
              {/each}
            </ol>
            {/each}
        {:else}
            No commits.
        {/if}

</div>