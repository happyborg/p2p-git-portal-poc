<script>
    export let dismissable;
    export let errorMessage;
    export let type
     let types = {
        error: 'text-red-600 bg-red-50 border-red-100',
        info: 'text-gray-600 bg-blue-50 border-blue-100'
     }
</script>
<div class="{types[type] ? types[type] : types['info']} my-5 block text-sm text-left border p-4 rounded-md relative" role="alert">
    <slot></slot>
    {#if dismissable && dismissable === "true"}
        <span class="absolute top-0 bottom-0 right-0 px-4 py-3" on:click={() => { errorMessage = undefined; }}>
            <svg class="fill-current h-6 w-6 text-red-600" role="button" xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 20 20">
                <title>Dismiss</title>
                <path
                    d="M14.348 14.849a1.2 1.2 0 0 1-1.697 0L10 11.819l-2.651 3.029a1.2 1.2 0 1 1-1.697-1.697l2.758-3.15-2.759-3.152a1.2 1.2 0 1 1 1.697-1.697L10 8.183l2.651-3.031a1.2 1.2 0 1 1 1.697 1.697l-2.758 3.152 2.758 3.15a1.2 1.2 0 0 1 0 1.698z" />
            </svg>
        </span>
    {/if}
</div>
