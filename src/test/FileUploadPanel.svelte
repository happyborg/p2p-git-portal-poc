<script>
import { onMount } from "svelte";

export let droppedFiles;

// Based on: https://stackoverflow.com/questions/3590058/does-html5-allow-drag-drop-upload-of-folders-or-a-folder-tree/53058574
// NOTE: directory upload works in Firefox, but not in Chromium 82, where
// readEntries() always returns an empty array.
//
// The docs suggest it should work in Chrome but say the API is deprecated anyway
// See: https://developer.mozilla.org/en-US/docs/Web/API/FileSystemDirectoryReader/readEntries

// Drop handler function to get all files
async function getAllFileEntries(dataTransferItemList) {
  let fileEntries = [];
  // Use BFS to traverse entire directory/file structure
  let queue = [];
  // Unfortunately dataTransferItemList is not iterable i.e. no forEach
  for (let i = 0; i < dataTransferItemList.length; i++) {
    queue.push(dataTransferItemList[i].webkitGetAsEntry());
  }
  while (queue.length > 0) {
    let entry = queue.shift();
    if (entry.isFile) {
    //   console.log('file:', entry.fullPath);
      fileEntries.push(entry);
    } else if (entry.isDirectory) {
    //   console.log('directory:', entry.fullPath);
      let reader = entry.createReader();
      queue.push(...await readAllDirectoryEntries(reader));
    }
  }
  return fileEntries;
}

// Get all the entries (files or sub-directories) in a directory by calling readEntries until it returns empty array
async function readAllDirectoryEntries(directoryReader) {
  let entries = [];
  let readEntries = await readEntriesPromise(directoryReader);
  while (readEntries.length > 0) {
    entries.push(...readEntries);
    readEntries = await readEntriesPromise(directoryReader);
  }
  return entries;
}

// Wrap readEntries in a promise to make working with readEntries easier
async function readEntriesPromise(directoryReader) {
  try {
    return await new Promise((resolve, reject) => {
      directoryReader.readEntries(resolve, reject);
    });
  } catch (err) {
    console.log(err);
  }
}

let elDrop;
// let elItems;

onMount(() => {
    elDrop = document.getElementById('dropzone');
    // elItems = document.getElementById('items');
    
    elDrop.addEventListener('dragover', function (event) {
        event.preventDefault();
        // elItems.innerHTML = 0;
    });
    
    elDrop.addEventListener('drop', async function (event) {
        // console.dir(event);
        event.preventDefault();
        let items = await getAllFileEntries(event.dataTransfer.items);
        droppedFiles = [...items];
        // elItems.innerHTML = droppedFiles.length;
    });
});
</script>

<style>
#dropzone {
    background-color: #cfc;
    border: solid 3px #9c9;
    color: #9c9;
    min-height: 50px;
    padding: 20px;
    text-shadow: 1px 1px 0 #fff;
}
</style>

<div>
    <hr>
    <h2>Upload Repository</h2>
    <p>Note: directory upload works in Firefox but not Chromium 82</p>
    <div id="dropzone" effectAllowed="move">Drop repository here</div>
    <slot></slot>
</div>

