<script lang="ts">
  let resultText = $state("Drop your slack import zip file here");
  let isDragging = $state(false);
  let dragCounter = $state(0);

  // --- Drag & Drop Event Handlers ---
  function handleDragEnter(event: DragEvent): void {
    event.preventDefault();
    dragCounter++;
    if (dragCounter === 1) {
      isDragging = true;
    }
  }

  function handleDragOver(event: DragEvent): void {
    event.preventDefault();
    if (!isDragging) {
      isDragging = true;
    }
  }

  function handleDragLeave(event: DragEvent): void {
    dragCounter--;
    if (dragCounter === 0) {
      isDragging = false;
    }
  }

  function handleDrop(event: DragEvent): void {
    event.preventDefault();
    dragCounter = 0;
    isDragging = false;
  }  
</script>

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
    <!-- Slack Logo SVG -->
    <svg xmlns="http://www.w3.org/2000/svg" x="0px" y="0px" width="100" height="100" viewBox="0 0 48 48">
      <path fill="#33d375" d="M33,8c0-2.209-1.791-4-4-4s-4,1.791-4,4c0,1.254,0,9.741,0,11c0,2.209,1.791,4,4,4s4-1.791,4-4	C33,17.741,33,9.254,33,8z"></path><path fill="#33d375" d="M43,19c0,2.209-1.791,4-4,4c-1.195,0-4,0-4,0s0-2.986,0-4c0-2.209,1.791-4,4-4S43,16.791,43,19z"></path><path fill="#40c4ff" d="M8,14c-2.209,0-4,1.791-4,4s1.791,4,4,4c1.254,0,9.741,0,11,0c2.209,0,4-1.791,4-4s-1.791-4-4-4	C17.741,14,9.254,14,8,14z"></path><path fill="#40c4ff" d="M19,4c2.209,0,4,1.791,4,4c0,1.195,0,4,0,4s-2.986,0-4,0c-2.209,0-4-1.791-4-4S16.791,4,19,4z"></path><path fill="#e91e63" d="M14,39.006C14,41.212,15.791,43,18,43s4-1.788,4-3.994c0-1.252,0-9.727,0-10.984	c0-2.206-1.791-3.994-4-3.994s-4,1.788-4,3.994C14,29.279,14,37.754,14,39.006z"></path><path fill="#e91e63" d="M4,28.022c0-2.206,1.791-3.994,4-3.994c1.195,0,4,0,4,0s0,2.981,0,3.994c0,2.206-1.791,3.994-4,3.994	S4,30.228,4,28.022z"></path><path fill="#ffc107" d="M39,33c2.209,0,4-1.791,4-4s-1.791-4-4-4c-1.254,0-9.741,0-11,0c-2.209,0-4,1.791-4,4s1.791,4,4,4	C29.258,33,37.746,33,39,33z"></path><path fill="#ffc107" d="M28,43c-2.209,0-4-1.791-4-4c0-1.195,0-4,0-4s2.986,0,4,0c2.209,0,4,1.791,4,4S30.209,43,28,43z"></path>
    </svg>
    <!-- File Icon (SVG) - Lucide-style for minimalism -->
    <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="file-icon">
        <path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z"/>
        <path d="M14 2v4a2 2 0 0 0 2 2h4"/>
        <path d="M10 9H8"/>
        <path d="M16 13H8"/>
        <path d="M16 17H8"/>
    </svg>
  </div>

  <div class="result-text">
    {resultText}
  </div>
</div>

<style>

  .drop-zone {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 2rem;
    border-radius: 0.75rem;
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
    transition: all 300ms ease-in-out;
    border: 2px dashed #4a5568;
    background-color: #2d3748;
    min-width: 300px;
    min-height: 200px;
  }

  .drop-zone.is-dragging {
    border-color: #3b82f6;
    /* Adjusted box-shadow for a more pronounced glow */
    box-shadow: 0 0 30px rgba(59, 130, 246, 0.7), 0 0 15px rgba(59, 130, 246, 0.5);
    transform: scale(1.05);
  }

  .icons-container {
    display: flex;
    align-items: center;
    margin-bottom: 1rem;
  }

  .file-icon {
    color: #60a5fa;
    /* No margin-bottom here, it's in the container */
  }

  

  .result-text {
    font-size: 1.25rem;
    font-weight: 500;
    text-align: center;
    margin-bottom: 1rem;
  }
</style>
