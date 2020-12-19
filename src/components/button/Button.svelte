<script>
    export let type
     let types = {
        white: 'bg-white text-gray-800 border-gray-200 hover:border-red-600 hover:bg-red-500 hover:text-white',
        red: 'bg-red-600 text-white border-gray-200 hover:border-red-800 hover:bg-red-500 hover:text-white'
     }
</script>

<button class="{types[type] ? types[type] : types['white']} text-sm font-bold border py-2 px-5 inline-flex items-center">
    <slot></slot>
</button>