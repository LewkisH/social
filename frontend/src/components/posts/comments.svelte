<script>
  import { selectUser } from "../../utils";
  import { IMAGE_URL } from "../../stores";
  import { fade, slide } from "svelte/transition";
  export let comment;
  $: showTime = false;
</script>

<!-- svelte-ignore a11y-mouse-events-have-key-events -->
<div class="singleComment">
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div class="userInfo" on:click={() => selectUser(comment.userID)}>
    <p class="commentCreator">
      {comment.user.firstName}
      {comment.user.lastName}
    </p>
    <p class="commentCreatorAvatar">
      <img
        class="postCreatorAvatar"
        src={IMAGE_URL}{comment.user.avatar}
        alt="user avatar"
      />
    </p>
  </div>

  <div
    class="commentContent"
    on:mouseover={() => {
      showTime = true;
    }}
    on:mouseout={() => {
      showTime = false;
    }}
  >
    {comment.content}
    {#if comment.img}
      <img src={IMAGE_URL}{comment.img} alt="This is a comment" />
    {/if}
  </div>
  {#if showTime === true}
    <p
      class="commentCreatedAt"
      transition:fade={{ delay: 1000, duration: 350 }}
    >
      {comment.createdAt}
    </p>
  {/if}
</div>

<style>
  .commentContent img {
    margin-top: 10px;
    border: 1px solid rgba(255, 255, 255, 0.447);
    border-radius: 10px;
  }
  div {
    padding: 8px;
    border-radius: 40px;
  }

  p {
    margin: 0;
  }

  .commentCreatedAt {
    margin: 0;
    padding: 0;
  }

  .postCreatorAvatar {
    padding: 4px 0;
    border-radius: 1000px;
    max-height: 40px;
  }
  .userInfo {
    cursor: pointer;
    width: 90px;
  }
  .singleComment {
    align-items: center;
    padding: 0;
    font-size: small;
    display: flex;
  }

  .commentContent {
    display: flex;
    border: solid 1px #333;
    justify-content: center;
    flex-direction: column;
    max-width: 60%;
    margin: 0 8px;
  }
  .commentContent > img {
    max-width: 300px;
    max-height: 300px;
  }
</style>
