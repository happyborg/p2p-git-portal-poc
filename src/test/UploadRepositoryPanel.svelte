<script>
import Button from "../components/button/Button.svelte";
export let uploadRoot;
export let filesToUpload;

let localPath = ''

function handleChange(event) {
  // console.log('handleChange: ', localPath);
  // console.dir(event);
  // console.dir(this);
  // console.dir(this.files)
  // let fileInfo = this.files[0];
  // console.dir(fileInfo);

  // uploadRepository(this.files);
}

function handleInput(event) {
  console.log('handleInput: ', localPath);
  // console.dir(event);
  // console.dir(this);
  // console.dir(this.files)
  let fileInfo = this.files[0];
  if (fileInfo !== undefined) {
    let root = ''
    root = fileInfo.webkitRelativePath.split('/')[0]
    // console.dir(fileInfo); 
    uploadRepository(root, this.files);
  }
}

function uploadRepository(root, files){
  console.log("uploadRepository() files.length = ", files.length)
  uploadRoot = root
  filesToUpload = [...files]
}

</script>

<style>
</style>

<div class="mt-6 pt-4 border-t">
    <label class="inline-block text-gray-800 text-m font-semibold mb-2" for="">Upload repository</label> 
    <div class="overflow-hidden relative w-64">
      <Button type="white"> 
        <span class="mr-2">Choose file</span>
        <svg fill="currentColor" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16"><path fill-rule="evenodd" d="M8.53 1.22a.75.75 0 00-1.06 0L3.72 4.97a.75.75 0 001.06 1.06l2.47-2.47v6.69a.75.75 0 001.5 0V3.56l2.47 2.47a.75.75 0 101.06-1.06L8.53 1.22zM3.75 13a.75.75 0 000 1.5h8.5a.75.75 0 000-1.5h-8.5z"></path></svg>
      </Button>
        <input
            class="cursor-pointer absolute block py-2 px-4 w-full opacity-0 pin-r pin-t"
            type="file"
            webkitdirectory=true bind:value={localPath} on:change={handleChange} on:input={handleInput}
        >
    </div>
     <slot></slot>
</div>

