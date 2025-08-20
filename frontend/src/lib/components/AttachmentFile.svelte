<script lang="ts">
let { file } = $props()
</script>

<a
    href={file.url_private}
    target="_blank"
    rel="noopener noreferrer"
    download={file.name}
    class="file-preview generic-file"
>
    <div class="file-icon"></div>
    <div class="file-info">
        <div class="file-name">
            {file.name || "file_not_found"}
        </div>
        <div class="file-type">{file.mimetype || "unknown"}</div>
    </div>
</a>

<style>
    /* Container for the file preview */
    .file-preview {
        display: flex;
        align-items: center;
        background-color: #f3f4f6;
        border: 1px solid #e1e1e1;
        border-radius: 8px;
        padding: 8px 12px;
        margin: 8px 0;
        color: #333;
        width: 300px;
        height: 50px;
        text-decoration: none;
    }

    .file-preview:hover {
        background-color: #e8e8e8; /* Slightly darker grey on hover */
    }

    /* Icon container on the left */
    .file-icon {
        width: 48px;
        height: 48px;
        background-color: #d1d5db;
        border-radius: 4px;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-right: 12px;
    }

    .file-icon::before {
        content: "V";
        font-size: 24px;
        font-weight: 700;
    }

    .file-info {
        display: flex;
        flex-direction: column;
        justify-content: center;
        overflow: hidden;
    }

    .file-name {
        font-size: 20px;
        font-weight: 600;
        color: #444;
        white-space: nowrap;
        text-overflow: ellipsis;
    }

    .file-type {
        font-size: 12px;
        color: #666;
        white-space: nowrap;
        overflow: hidden;
    }
</style>
