<script>
// Select active repository from those stored in the Go filesystem

export let activeRepository;
export let allRepositories;

</script>

<style>
.highlight {
    font-weight: bold;
}
.highlight svg {
    color: #f0572a;
}

</style>
    <h2
    class="text-gray-800 text-m font-semibold mb-2 mt-8"
    >
    Repositories</h2>
    <div class="w-full border-solid	overflow-hidden shadow-xs mt-2 py-2 border-t border-gray-300 overflow-hidden text-sm font-medium text-gray-500">
    <div class="w-full overflow-x-auto">
    
             <ul class="">
                {#if allRepositories !== undefined && allRepositories.length > 0}
                    {#each allRepositories as repo, index}       
                    <li class:highlight={index == activeRepository} class="focus:bg-gray-50 hover:bg-gray-transDark flex my-2 items-center">
                        <svg xmlns="http://www.w3.org/2000/svg" 
                                class="w-4 h-4 mr-2"
                                fill="currentColor"
                                viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M2 2.5A2.5 2.5 0 014.5 0h8.75a.75.75 0 01.75.75v12.5a.75.75 0 01-.75.75h-2.5a.75.75 0 110-1.5h1.75v-2h-8a1 1 0 00-.714 1.7.75.75 0 01-1.072 1.05A2.495 2.495 0 012 11.5v-9zm10.5-1V9h-8c-.356 0-.694.074-1 .208V2.5a1 1 0 011-1h8zM5 12.25v3.25a.25.25 0 00.4.2l1.45-1.087a.25.25 0 01.3 0L8.6 15.7a.25.25 0 00.4-.2v-3.25a.25.25 0 00-.25-.25h-3.5a.25.25 0 00-.25.25z"></path>
                        </svg>
                         <div class="flex flex-col p-0"><div><a href= '#none' on:click={() => activeRepository = index}>{repo.path}</a></div> </div>
                    </li>
                    {/each}
                    
                {:else}
                <li class="text-xs">
                        <strong>No repositories found.</strong><br/>
                        Add repositories using New, Clone or Upload.
                </li>
                {/if}
    
    
                    </ul>
          </div>
    
</div>