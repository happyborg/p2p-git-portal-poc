<script>

export let width = "600px"

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

<div>
    <h3>Upload</h3>
    <p>
        <input width={width} type="file" webkitdirectory=true bind:value={localPath} on:change={handleChange} on:input={handleInput} />
    </p>
    <slot></slot>
</div>

