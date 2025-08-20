<script lang="ts">
// Svelte 5 state management with runes
let resultText = $state('Drop your Slack import zip file here')
let isDragging = $state(false)
let dragCounter = $state(0)

// --- Drag & Drop Event Handlers ---
// These are standard DOM event handlers and work perfectly with Svelte 5
function handleDragEnter(event: DragEvent): void {
  event.preventDefault()
  dragCounter++
  if (dragCounter === 1) {
    isDragging = true
  }
}

function handleDragOver(event: DragEvent): void {
  event.preventDefault()
  // This ensures the dragging state is active even if the enter event was missed
  if (!isDragging) {
    isDragging = true
  }
}

function handleDragLeave(event: DragEvent): void {
  event.preventDefault()
  dragCounter--
  if (dragCounter === 0) {
    isDragging = false
  }
}

function handleDrop(event: DragEvent): void {
  event.preventDefault()
  dragCounter = 0
  isDragging = false

  // You would typically process the dropped files here
  const files = event.dataTransfer?.files
  if (files && files.length > 0) {
    resultText = `${files.length} file(s) selected: ${files[0].name}`
    console.log(files)
  } else {
    resultText = 'No files were dropped.'
  }
}
</script>

<main>
  <div
    class="drop-zone"
    class:is-dragging={isDragging}
    role="button"
    tabindex="0"
    ondragenter={handleDragEnter}
    ondragleave={handleDragLeave}
    ondragover={handleDragOver}
    ondrop={handleDrop}
  >
    <div class="icons-container">
      <svg xmlns="http://www.w3.org/2000/svg" x="0px" y="0px" width="100" height="100" viewBox="0 0 48 48">
        <path fill="#33d375" d="M33,8c0-2.209-1.791-4-4-4s-4,1.791-4,4c0,1.254,0,9.741,0,11c0,2.209,1.791,4,4,4s4-1.791,4-4 C33,17.741,33,9.254,33,8z"></path><path fill="#33d375" d="M43,19c0,2.209-1.791,4-4,4c-1.195,0-4,0-4,0s0-2.986,0-4c0-2.209,1.791-4,4-4S43,16.791,43,19z"></path><path fill="#40c4ff" d="M8,14c-2.209,0-4,1.791-4,4s1.791,4,4,4c1.254,0,9.741,0,11,0c2.209,0,4-1.791,4-4s-1.791-4-4-4 C17.741,14,9.254,14,8,14z"></path><path fill="#40c4ff" d="M19,4c2.209,0,4,1.791,4,4c0,1.195,0,4,0,4s-2.986,0-4,0c-2.209,0-4-1.791-4-4S16.791,4,19,4z"></path><path fill="#e91e63" d="M14,39.006C14,41.212,15.791,43,18,43s4-1.788,4-3.994c0-1.252,0-9.727,0-10.984  c0-2.206-1.791-3.994-4-3.994s-4,1.788-4,3.994C14,29.279,14,37.754,14,39.006z"></path><path fill="#e91e63" d="M4,28.022c0-2.206,1.791-3.994,4-3.994c1.195,0,4,0,4,0s0,2.981,0,3.994c0,2.206-1.791,3.994-4,3.994  S4,30.228,4,28.022z"></path><path fill="#ffc107" d="M39,33c2.209,0,4-1.791,4-4s-1.791-4-4-4c-1.254,0-9.741,0-11,0c-2.209,0-4,1.791-4,4s1.791,4,4,4  C29.258,33,37.746,33,39,33z"></path><path fill="#ffc107" d="M28,43c-2.209,0-4-1.791-4-4c0-1.195,0-4,0-4s2.986,0,4,0c2.209,0,4,1.791,4,4S30.209,43,28,43z"></path>
      </svg>
    </div>
    <div class="result-text">
      {resultText}
    </div>
  </div>
</main>

<style>
 
  main {
    margin: 0;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
    color: #e2e8f0; /* Light gray text for better contrast */

    display: grid;
    place-items: center;
    min-height: 100vh;
    background-color: #4a0e4e; /* Aubergine color */
  }

  .drop-zone {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 2rem;
    border-radius: 0.75rem;
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
    transition: all 300ms ease-in-out;
    border: 2px dashed #6b7280; /* A slightly lighter gray for the dash */
    background-color: #2d3748;
    min-width: 350px;
    min-height: 220px;
    cursor: pointer;
  }

  .drop-zone.is-dragging {
    border-color: #3b82f6;
    box-shadow: 0 0 30px rgba(59, 130, 246, 0.7), 0 0 15px rgba(59, 130, 246, 0.5);
    transform: scale(1.05);
  }

  .icons-container {
    display: flex;
    align-items: center;
    gap: 1rem; /* Adds some space between the icons */
    margin-bottom: 1rem;
  }

  .file-icon {
    color: #60a5fa;
  }

  .result-text {
    font-size: 1.1rem; /* Slightly adjusted for better fit */
    font-weight: 500;
    text-align: center;
    margin-top: 0.5rem; /* Replaced margin-bottom with margin-top for better spacing */
  }
</style>