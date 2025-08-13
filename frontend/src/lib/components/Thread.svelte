<script lang="ts">
    import { LogError, LogInfo } from "../../../wailsjs/runtime";
    import { GetUserAvatar } from "../../../wailsjs/go/app/App";

    const { replyUsers, replyCount } = $props();
    const cache = $state(new Map<string, string>());
    let avatars: string[] = $state([]);

    async function getUser(userId: string): Promise<string> {
        if (cache.has(userId)) {
            LogInfo(`Cache hit for user ${userId}`);
            return cache.get(userId)!;
        }

        try {
            const url = await GetUserAvatar(userId);
            cache.set(userId, url);

            return url;
        } catch (error) {
            LogError(`Failed to get user info for ${userId}:`);
            return "";
        }
    }

    async function getDedupedAvatars(users: string[]): Promise<string[]> {
        const dedupedUsers = Array.from(new Set(users));

        const avatars = await Promise.all(
            dedupedUsers.map((userId) => getUser(userId)),
        );

        // add a slice here if want to trim the list
        return avatars;
    }

    $effect(() => {
        (async () => {
            avatars = await getDedupedAvatars(replyUsers);
        })();
    });
</script>

<div class="thread-info">
    <div class="reply-avatars">
        {#each avatars as avatar_url}
            <img class="reply-avatar" src={avatar_url} alt="Reply user" />
        {/each}
    </div>

    <a href="#" class="reply-link">
        {replyCount}
        {replyCount === 1 ? "reply" : "replies"}
    </a>

    <span class="view-thread-prompt">View thread</span>
</div>

<style>
    .thread-info {
        display: flex;
        align-items: center;
        margin-top: 0.5rem;
        cursor: pointer; /* Make the whole area feel clickable */
    }

    .reply-avatars {
        display: flex;
    }

    .reply-avatar {
        width: 24px;
        height: 24px;
        border-radius: 50%;
        border: 2px solid white; /* Creates separation between avatars */
        object-fit: cover;
    }

    /* This selector makes the avatars overlap */
    .reply-avatar:not(:first-child) {
        margin-left: -8px;
    }

    .reply-link {
        color: #1264a3; /* A nice blue, similar to Slack's */
        font-weight: 600;
        font-size: 0.9rem;
        text-decoration: none;
        margin-left: 0.5rem;
    }

    .reply-link:hover {
        text-decoration: underline;
    }

    .view-thread-prompt {
        color: #9e9ea6;
        font-size: 0.9rem;
        margin-left: 0.5rem;
    }
</style>
